package browser

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/chromedp/cdproto/cdp"
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

// projectDesc contains a url, description for a project.
type projectDesc struct {
	URL, Description string
}

// ListAwesomeGoProjects is the highest level logic for browsing to the
// awesome-go page, finding the specified section sect, and retrieving the
// associated projects from the page.
func ListAwesomeGoProjects(ctx context.Context, URL string, sect string) (map[string]projectDesc, error) {
	sel := fmt.Sprintf(`//p[text()[contains(., '%s')]]`, sect)

	if err := chromedp.Run(ctx, chromedp.Navigate(URL)); err != nil {
		return nil, fmt.Errorf("could not navigate to github: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitVisible(sel)); err != nil {
		return nil, fmt.Errorf("could not get section: %v", err)
	}

	sib := sel + `/following-sibling::ul/li`
	// get project link text
	var projects []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/child::a/text()`, &projects)); err != nil {
		return nil, fmt.Errorf("could not get projects: %v", err)
	}

	// get links and description text
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/child::node()`, &linksAndDescriptions)); err != nil {
		return nil, fmt.Errorf("could not get links and descriptions: %v", err)
	}

	// check length
	if 2*len(projects) != len(linksAndDescriptions) {
		return nil, fmt.Errorf("projects and links and descriptions lengths do not match (2*%d != %d)", len(projects), len(linksAndDescriptions))
	}

	// process data
	res := make(map[string]projectDesc)
	for i := 0; i < len(projects); i++ {
		res[projects[i].NodeValue] = projectDesc{
			URL:         linksAndDescriptions[2*i].AttributeValue("href"),
			Description: strings.TrimPrefix(strings.TrimSpace(linksAndDescriptions[2*i+1].NodeValue), "- "),
		}
	}

	return res, nil
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
