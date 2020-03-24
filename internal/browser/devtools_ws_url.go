package browser

import (
	"encoding/json"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

const (
	devtoolsWsScheme          = "ws://"
	devtoolsWsDomain          = "headless-browser:9222"
	devtoolsWsJSONVersionPath = "/json/version"
	devtoolsEndpointPath      = "/devtools/browser/"
)

// GetDevtoolsEndpoint gets the string dev tools endpoint.
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
