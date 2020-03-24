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
	devtoolsWs, error := getDevtoolsWs()

	if error != nil {
		return "", error
	}

	wsDebuggerURL, error := getWsDebuggerURL(devtoolsWs)

	if error != nil {
		return "", error
	}

	return wsDebuggerURL, nil
}

func getWsDebuggerURL(devtoolsWs []byte) (string, error) {
	var devtoolsWsJSON map[string]interface{}

	if error := json.Unmarshal([]byte(devtoolsWs), &devtoolsWsJSON); error != nil {
		return "", error
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
	out, error := cmd.Output()

	if error != nil {
		return []byte(""), error
	}

	return out, nil
}
