// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/ph-piment/headless-chrome/code/go/browser"

	"github.com/chromedp/chromedp"
)

func main() {
	ctxt, allocCancel, ctxtCancel := getContext()
	defer allocCancel()
	defer ctxtCancel()

	// run task list
	var body string
	if err := chromedp.Run(ctxt,
		chromedp.Navigate("https://duckduckgo.com"),
		chromedp.WaitVisible("#logo_homepage_link"),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		log.Fatalf("Failed getting body of duckduckgo.com: %v", err)
	}

	log.Println("Body of duckduckgo.com starts with:")
	log.Println(body[0:100])
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
