package redis

import (
	"reflect"
	"testing"
)

func TestHSetSuccess(t *testing.T) {
	expected := "ccc"
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	field := "bbb"
	setErr := HSet(key, field, "ccc")
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := HGet(key, field)
	if actual != expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if getErr != nil {
		t.Errorf("got error: %v\n", getErr)
	}
	afterDelErr := DelString(key)
	if afterDelErr != nil {
		t.Errorf("got error: %v\n", afterDelErr)
	}
}

func TestHGetAllSuccess(t *testing.T) {
	expected := "ccc"
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := HSet(key, "bbb", "ccc")
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	setErr = HSet(key, "ddd", "eee")
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := HGetAll(key)
	if !reflect.DeepEqual(actual, map[string]string{"bbb": "ccc", "ddd": "eee"}) {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if getErr != nil {
		t.Errorf("got error: %v\n", getErr)
	}
	afterDelErr := DelString(key)
	if afterDelErr != nil {
		t.Errorf("got error: %v\n", afterDelErr)
	}
}
