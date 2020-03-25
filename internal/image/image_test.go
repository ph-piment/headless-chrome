package image

import (
	"testing"
)

var aaaTextFilePath = "/tmp/aaa"
var notExistsFilePath = "/tmpaaa/bbb"

func TestWriteImageByByteSuccess(t *testing.T) {
	actual := WriteImageByByte([]byte("aaa"), aaaTextFilePath)
	var expected error
	expected = nil
	if actual != expected {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
}

func TestWriteImageByByteFailed(t *testing.T) {
	actual := WriteImageByByte([]byte("bbb"), notExistsFilePath)
	var expected error
	expected = nil
	if actual == nil {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
}

/*
func TestReadImageByPathSuccess(t *testing.T) {
	actual, err := ReadImageByPath("/go/src/work/outputs/images/compare/source/image.png")
	var expected error
	expected = nil
	if actual == nil {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if err != nil {
		t.Errorf("got error: %v\n", err)
	}
}
*/

func TestReadImageByPathFailed(t *testing.T) {
	actual, err := ReadImageByPath(notExistsFilePath)
	var expected error
	expected = nil
	if actual != nil {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if err == nil {
		t.Errorf("got error: %v\n", err)
	}

	actual, err = ReadImageByPath(aaaTextFilePath)
	expected = nil
	if actual != nil {
		t.Errorf("got: %v\nwant: %v\n", actual, expected)
	}
	if err == nil {
		t.Errorf("got error: %v\n", err)
	}
}
