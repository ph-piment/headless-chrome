package redis

import (
	"reflect"
	"testing"
)

func TestLPushSuccess(t *testing.T) {
	expected := []string{"val3", "val2", "val1"}
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := LPush(key, []string{"val1", "val2", "val3"})
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := AllRange(key)
	if !reflect.DeepEqual(actual, expected) {
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

func TestRPushSuccess(t *testing.T) {
	expected := []string{"val1", "val2", "val3"}
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := RPush(key, expected)
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := AllRange(key)
	if !reflect.DeepEqual(actual, expected) {
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

func TestLSetSuccess(t *testing.T) {
	expected := []string{"val1", "val4"}
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := RPush(key, []string{"val1", "val2", "val3"})
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	lSetErr := LSet(key, 1, "val4")
	if lSetErr != nil {
		t.Errorf("got error: %v\n", lSetErr)
	}
	val, lIndexErr := LIndex(key, 1)
	if val != "val4" {
		t.Errorf("got: %v\nwant: %v\n", val, "val4")
	}
	if lIndexErr != nil {
		t.Errorf("got error: %v\n", lIndexErr)
	}
	actual, getErr := LRange(key, 0, 1)
	if !reflect.DeepEqual(actual, expected) {
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
