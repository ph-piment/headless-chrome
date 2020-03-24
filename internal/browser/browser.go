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

// GetContext get context by NewRemoteAllocator.
func GetContext() (context.Context, context.CancelFunc, context.CancelFunc) {
	devtoolsEndpoint, error := GetDevtoolsEndpoint()
	if error != nil {
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
			_, _, contentSize, error := page.GetLayoutMetrics().Do(ctx)
			if error != nil {
				return error
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			error = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
			if error != nil {
				return error
			}

			*res, error = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if error != nil {
				return error
			}
			return nil
		}),
	}
}

func openImage(path string) (image.Image, error) {
	f, error := os.Open(path)
	if error != nil {
		return nil, errors.Wrap(error, "failed to open")
	}
	defer f.Close()

	img, _, error := image.Decode(f)
	if error != nil {
		return nil, errors.Wrap(error, "failed to decode image")
	}
	return img, nil
}

func GetImageByURL(ctx context.Context, url string, imagePath string) image.Image {
	var buf []byte
	if error := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); error != nil {
		log.Fatal(error)
	}
	if error := ioutil.WriteFile(imagePath, buf, 0644); error != nil {
		log.Fatal(error)
	}
	imagefile, error := openImage(imagePath)
	if error != nil {
		log.Fatal(error)
	}

	return imagefile
}
