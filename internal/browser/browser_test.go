package browser

import (
	"reflect"
	"testing"
)

func TestGetFullScreenshotSuccess(t *testing.T) {
	var buf []byte

	tasks := getFullScreenshot("https://www.google.com", 90, &buf)
	actual := reflect.TypeOf(tasks).String()
	expected := "chromedp.Tasks"
	if actual != expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
}
