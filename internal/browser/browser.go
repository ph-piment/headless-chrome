package browser

import (
	"context"
	"log"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ScreenshotQuality is screenshot quality
const ScreenshotQuality = 50

// GetContext get context by NewRemoteAllocator.
func GetContext() (context.Context, context.CancelFunc, context.CancelFunc) {
	devtoolsEndpoint, err := GetDevtoolsEndpoint()
	if err != nil {
		log.Fatal("must get devtools endpoint")
	}

	allocCtx, allocCxl :=
		chromedp.NewRemoteAllocator(context.Background(), devtoolsEndpoint)

	ctx, ctxCxl := chromedp.NewContext(allocCtx)

	return ctx, allocCxl, ctxCxl
}

func getFullScreenshot(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height :=
				int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
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

// GetFullScreenshotByteByURL get image by URL.
func GetFullScreenshotByteByURL(ctx context.Context, url string) ([]byte, error) {
	var buf []byte
	if err := chromedp.Run(ctx, getFullScreenshot(url, ScreenshotQuality, &buf)); err != nil {
		return nil, err
	}
	return buf, nil
}
