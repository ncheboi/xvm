package gzip_test

import (
	"bytes"
	"testing"

	"."
)

func TestGzipGunzip(t *testing.T) {
	expected := []byte("test")
	actual, err := gzip.Decompress(gzip.Compress(expected))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
