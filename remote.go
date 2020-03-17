// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"log"
	"os/exec"
	"encoding/json"
	"io"
	"strings"

	"github.com/chromedp/chromedp"
)

func getDebugURL() string {
	ws := "ws://"
	domain := "headless-browser:9222"
	jsonPath := domain + "/json/version"
	cmd := exec.Command("curl", "-H", "host:", jsonPath)
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, "hoge")
	stdin.Close()
	out, _ := cmd.Output()

	var result map[string]interface{}

	if err := json.Unmarshal([]byte(out), &result); err != nil {
		log.Fatal(err)
	}
	return ws + domain + strings.Replace(result["webSocketDebuggerUrl"].(string), ws, "", 1)
}

var flagDevToolWsUrl = flag.String("devtools-ws-url", getDebugURL(), "DevTools WebSsocket URL")

func main() {
	flag.Parse()

	if *flagDevToolWsUrl == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	// create allocator context for use with creating a browser context later
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), *flagDevToolWsUrl)
	defer cancel()

	// create context
	ctxt, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

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
