package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(src []byte) []byte {
	dst := new(bytes.Buffer)

	writer := gzip.NewWriter(dst)
	writer.Write(src)
	writer.Close()

	return dst.Bytes()
}

func Decompress(src []byte) ([]byte, error) {
	srcBuf := bytes.NewBuffer(src)
	dstBuf := new(bytes.Buffer)

	reader, err := gzip.NewReader(srcBuf)
	if err == nil {
		_, err = io.Copy(dstBuf, reader)
	}
	return dstBuf.Bytes(), err
}
