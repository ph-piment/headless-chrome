// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"image"
	"image/color"
	"image/png"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ph-piment/headless-chrome/code/go/browser"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/orisano/pixelmatch"
)

type colorValue color.RGBA

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage of pixelmatch [flags] image1 image2 :")
	}
	img1 := args[0]
	img2 := args[1]

	ctx, allocCancel, ctxtCancel := getContext()
	defer allocCancel()
	defer ctxtCancel()

	compareDir := filepath.Dir("/go/src/work/outputs/images/compare/")
	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(img1, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	sourceImagePath := compareDir + "/source/image.png"
	if err := ioutil.WriteFile(sourceImagePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	img1file, err := openImage(sourceImagePath)
	if err != nil {
		log.Fatal(err)
	}

	buf = nil
	if err := chromedp.Run(ctx, fullScreenshot(img2, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	targetImagePath := compareDir + "/target/image.png"
	if err := ioutil.WriteFile(targetImagePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	img2file, err := openImage(targetImagePath)
	if err != nil {
		log.Fatal(err)
	}

	// compare
	threshold := flag.Float64("threshold", 0.1, "threshold")
	aa := flag.Bool("aa", false, "ignore anti alias pixel")
	alpha := flag.Float64("alpha", 0.1, "alpha")
	antiAliased := colorValue(color.RGBA{R: 255, G: 255})
	diffColor := colorValue(color.RGBA{R: 255})
	var out image.Image
	opts := []pixelmatch.MatchOption{
		pixelmatch.Threshold(*threshold),
		pixelmatch.Alpha(*alpha),
		pixelmatch.AntiAliasedColor(color.RGBA(antiAliased)),
		pixelmatch.DiffColor(color.RGBA(diffColor)),
		pixelmatch.WriteTo(&out),
	}
	if *aa {
		opts = append(opts, pixelmatch.IncludeAntiAlias)
	}

	_, err = pixelmatch.MatchPixel(img1file, img2file, opts...)
	if err != nil {
		log.Fatal(err)
	}

	var w io.Writer
	f, err := os.Create(compareDir + "/result/image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w = f

	var encErr error
	encErr = png.Encode(w, out)
	if encErr != nil {
		log.Fatal(err)
	}
}

// TODO: move ...
func getContext() (context.Context, context.CancelFunc, context.CancelFunc) {
	devToolWsURL := getDevToolWsURL()

	// create allocator context for use with creating a browser context later
	allocatorContext, allocCancel := chromedp.NewRemoteAllocator(context.Background(), devToolWsURL)

	// create context
	ctxt, ctxtCancel := chromedp.NewContext(allocatorContext)

	return ctxt, allocCancel, ctxtCancel
}

// TODO: move browser package.
func getDevToolWsURL() string {
	flag.Parse()

	devtoolsEndpoint, err := browser.GetDevtoolsEndpoint()
	if err != nil {
		log.Fatal("must get devtools endpoint")
	}

	flagDevToolWsURL := flag.String("devtools-ws-url", devtoolsEndpoint, "DevTools WebSsocket URL")
	if *flagDevToolWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	return *flagDevToolWsURL
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open")
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}
	return img, nil
}
