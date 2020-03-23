package browser

import (
	"context"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/pkg/errors"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GetContext() (context.Context, context.CancelFunc, context.CancelFunc) {
	devtoolsEndpoint, err := GetDevtoolsEndpoint()
	if err != nil {
		log.Fatal("must get devtools endpoint")
	}

	allocatorContext, allocCancel := chromedp.NewRemoteAllocator(context.Background(), devtoolsEndpoint)

	ctxt, ctxtCancel := chromedp.NewContext(allocatorContext)

	return ctxt, allocCancel, ctxtCancel
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

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
	if err := ioutil.WriteFile(imagePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	imagefile, err := openImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	return imagefile
}
