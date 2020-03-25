package browser

import (
	"reflect"
	"testing"
)

/*
func TestGetDevtoolsEndpointSuccess(t *testing.T) {
	actual, err := GetDevtoolsEndpoint()
	expected :=
		"/^" + devtoolsWsScheme + devtoolsWsDomain + devtoolsEndpointPath + ".*$/"
	if regexp.MustCompile(expected).Match([]byte(actual)) {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if err != nil {
		t.Errorf("got error: %v\n", err)
	}
}
*/

func TestGetDevtoolsWsByteSuccess(t *testing.T) {
	devtoolsWsByte := getDevtoolsWsByte()
	actual := reflect.TypeOf(devtoolsWsByte).String()
	expected := "[]uint8"
	if actual != expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
}

func TestGetWsDebuggerURLSuccess(t *testing.T) {
	devtoolsWs := []byte(`
	{
   "Browser": "Chrome/80.0.3987.132",
   "Protocol-Version": "1.3",
   "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
   "V8-Version": "8.0.426.26",
   "WebKit-Version": "537.36 (@fcea73228632975e052eb90fcf6cd1752d3b42b4)",
   "webSocketDebuggerUrl": "ws:///devtools/browser/a10e17fa-7480-4225-88f2-3c824d6f9c88"
}`)

	wsDebuggerURL, err := getWsDebuggerURL(devtoolsWs)
	expected := "ws://headless-browser:9222/devtools/browser/a10e17fa-7480-4225-88f2-3c824d6f9c88"
	if wsDebuggerURL != expected {
		t.Errorf("got: %v\nwant: %v\n", wsDebuggerURL, expected)
	}
	if err != nil {
		t.Errorf("got error: %v\n", err)
	}

	devtoolsWs = []byte(`
	{
   "Browser": "Chrome/80.0.3987.132",
   "Protocol-Version": "1.3",
   "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
   "V8-Version": "8.0.426.26",
   "WebKit-Version": "537.36 (@fcea73228632975e052eb90fcf6cd1752d3b42b4)",
   "webSocketDebuggerUrl": "ws://headless-browser:9222/devtools/browser/a10e17fa-7480-4225-88f2-3c824d6f9c88"
}`)

	wsDebuggerURL, err = getWsDebuggerURL(devtoolsWs)
	if wsDebuggerURL != expected {
		t.Errorf("got: %v\nwant: %v\n", wsDebuggerURL, expected)
	}
	if err != nil {
		t.Errorf("got error: %v\n", err)
	}
}

func TestGetWsDebuggerURLFailed(t *testing.T) {
	devtoolsWs := []byte(`
	{
   "Browser": "Chrome/80.0.3987.132",
   "Protocol-Version": "1.3",
   "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
   "V8-Version": "8.0.426.26",
   "WebKit-Version": "537.36 (@fcea73228632975e052eb90fcf6cd1752d3b42b4)",
   "webSocketDebuggerUrl": "ws:///devtools/browser/a10e17fa-7480-4225-88f2-3c824d6f9c88"
`)

	wsDebuggerURL, err := getWsDebuggerURL(devtoolsWs)
	expected := ""
	if wsDebuggerURL != expected {
		t.Errorf("got: %v\nwant: %v\n", wsDebuggerURL, expected)
	}
	if err == nil {
		t.Errorf("got error: %v\n", err)
	}

	devtoolsWs = []byte(`
	{
   "Browser": "Chrome/80.0.3987.132",
   "Protocol-Version": "1.3",
   "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
   "V8-Version": "8.0.426.26",
   "WebKit-Version": "537.36 (@fcea73228632975e052eb90fcf6cd1752d3b42b4)",
   "ddd": "ws://headless-browser:9222/devtools/browser/a10e17fa-7480-4225-88f2-3c824d6f9c88"
}`)

	wsDebuggerURL, err = getWsDebuggerURL(devtoolsWs)
	if wsDebuggerURL != expected {
		t.Errorf("got: %v\nwant: %v\n", wsDebuggerURL, expected)
	}
	if err == nil {
		t.Errorf("got error: %v\n", err)
	}
}
