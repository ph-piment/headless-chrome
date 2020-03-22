package browser

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"math"
	"os/exec"
	"regexp"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	devtoolsWsScheme          = "ws://"
	devtoolsWsDomain          = "headless-browser:9222"
	devtoolsWsJSONVersionPath = "/json/version"
	devtoolsEndpointPath      = "/devtools/browser/"
)

func GetDevtoolsEndpoint() (string, error) {
	devtoolsWs, err := getDevtoolsWs()

	if err != nil {
		return "", err
	}

	wsDebuggerURL, err := getWsDebuggerURL(devtoolsWs)

	if err != nil {
		return "", err
	}

	return wsDebuggerURL, nil
}

func getWsDebuggerURL(devtoolsWs []byte) (string, error) {
	var devtoolsWsJSON map[string]interface{}

	if err := json.Unmarshal([]byte(devtoolsWs), &devtoolsWsJSON); err != nil {
		return "", err
	}

	webSocketDebuggerURL, ok := devtoolsWsJSON["webSocketDebuggerUrl"]
	if !ok {
		return "", errors.New("not exists webSocketDebuggerUrl key")
	}

	rep := regexp.MustCompile(devtoolsWsScheme + ".*" + devtoolsEndpointPath)
	devtoolsWsHash := rep.ReplaceAllString(webSocketDebuggerURL.(string), "")

	wsDebuggerURL := devtoolsWsScheme + devtoolsWsDomain + devtoolsEndpointPath + devtoolsWsHash
	return wsDebuggerURL, nil
}

func getDevtoolsWs() ([]byte, error) {
	// TODO: modify Http.get.(Http.get could not respond)
	devtoolsWsDomainJSONVersionPath := devtoolsWsDomain + devtoolsWsJSONVersionPath
	cmd := exec.Command("curl", "-H", "host:", devtoolsWsDomainJSONVersionPath)
	out, err := cmd.Output()

	if err != nil {
		return []byte(""), err
	}

	return out, nil
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

// TODO: move ...
func GetContext() (context.Context, context.CancelFunc, context.CancelFunc) {
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

	devtoolsEndpoint, err := GetDevtoolsEndpoint()
	if err != nil {
		log.Fatal("must get devtools endpoint")
	}

	flagDevToolWsURL := flag.String("devtools-ws-url", devtoolsEndpoint, "DevTools WebSsocket URL")
	if *flagDevToolWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	return *flagDevToolWsURL
}
