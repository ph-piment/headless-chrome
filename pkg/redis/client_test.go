package redis

import (
	"reflect"
	"testing"
)

func TestGetClientSuccess(t *testing.T) {
	actual := reflect.TypeOf(GetClient()).String()
	expected := "*redis.Client"
	if actual != expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
}
