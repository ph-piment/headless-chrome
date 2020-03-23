package browser

import (
	"context"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/orisano/pixelmatch"
)

type colorValue color.RGBA

var compareDir = filepath.Dir("/go/src/work/outputs/images/compare/")

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

func GetImageByURL(ctx context.Context, url string, imagePath string) image.Image {
	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	sourceImagePath := compareDir + imagePath
	if err := ioutil.WriteFile(sourceImagePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	imgfile, err := openImage(sourceImagePath)
	if err != nil {
		log.Fatal(err)
	}

	return imgfile
}

func DiffImage(sourceImage image.Image, targetImage image.Image, imagePath string) {
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

	_, err := pixelmatch.MatchPixel(sourceImage, targetImage, opts...)
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
