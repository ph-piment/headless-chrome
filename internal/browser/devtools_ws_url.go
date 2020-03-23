package browser

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"

	"github.com/chromedp/chromedp"
)

const (
	devtoolsWsScheme          = "ws://"
	devtoolsWsDomain          = "headless-browser:9222"
	devtoolsWsJSONVersionPath = "/json/version"
	devtoolsEndpointPath      = "/devtools/browser/"
)

func getDevtoolsEndpoint() (string, error) {
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

	devtoolsEndpoint, err := getDevtoolsEndpoint()
	if err != nil {
		log.Fatal("must get devtools endpoint")
	}

	flagDevToolWsURL := flag.String("devtools-ws-url", devtoolsEndpoint, "DevTools WebSsocket URL")
	if *flagDevToolWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	return *flagDevToolWsURL
}
