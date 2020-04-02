package redis

import (
	"testing"
	"time"
)

func TestSetStringSuccess(t *testing.T) {
	expected := "bbb"
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := SetString(key, expected)
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := GetString(key)
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

func TestSetStringWithExpireSuccess(t *testing.T) {
	expected := "bbb"
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	setErr := SetStringWithExpire(key, expected, time.Duration(1)*time.Second)
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	actual, getErr := GetString(key)
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

func TestSetStringWithExpireFailed(t *testing.T) {
	expected := "bbb"
	key := "aaa"
	delErr := DelString(key)
	if delErr != nil {
		t.Errorf("got error: %v\n", delErr)
	}
	waitSecond := 1
	setErr := SetStringWithExpire(key, expected, time.Duration(waitSecond)*time.Second)
	if setErr != nil {
		t.Errorf("got error: %v\n", setErr)
	}
	time.Sleep(time.Duration(waitSecond+1) * time.Second)
	actual, getErr := GetString(key)
	if actual == expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if getErr == nil {
		t.Errorf("got error: %v\n", getErr)
	}
	afterDelErr := DelString(key)
	if afterDelErr != nil {
		t.Errorf("got error: %v\n", afterDelErr)
	}
}
