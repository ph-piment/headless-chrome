// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"math"

	"github.com/ph-piment/headless-chrome/code/go/browser"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, allocCancel, ctxtCancel := getContext()
	defer allocCancel()
	defer ctxtCancel()

	// run task list
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://duckduckgo.com"),
		chromedp.WaitVisible("#logo_homepage_link"),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		log.Fatalf("Failed getting body of duckduckgo.com: %v", err)
	}

	log.Println("Body of duckduckgo.com starts with:")
	log.Println(body[0:100])

	// capture screenshot of an element
	CaptureScreenshotList := []map[string]string{
		{
			"url": "https://www.google.com/",
			"image": "google.png",
		},
		{
			"url": "https://www.yahoo.com/",
			"image": "yahoo.png",
		},
	}

	var buf []byte
	for index, CaptureScreenshot := range CaptureScreenshotList {
		log.Println("start index:", index)
		if err := chromedp.Run(ctx, fullScreenshot(CaptureScreenshot["url"], 90, &buf)); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(CaptureScreenshot["image"], buf, 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("done index:", index)
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
