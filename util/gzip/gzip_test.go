package gzip_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/skotchpine/xvm/util/gzip"
)

func TestGzipGunzip(t *testing.T) {
	expected := []byte("test")

	cxBuf := bytes.NewBuffer(expected)
	cxReader, _, err := gzip.Compress(cxBuf)
	if err != nil {
		t.Error(err)
	}

	dxReader, err := gzip.Decompress(cxReader)
	if err != nil {
		t.Error(err)
	}

	dxBuf := new(bytes.Buffer)
	if _, err := io.Copy(dxBuf, dxReader); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(dxBuf.Bytes(), expected) {
		t.Errorf("Expected %s, got %s", expected, dxBuf.Bytes())
	}
}
