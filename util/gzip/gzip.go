// Package gzip implements golang's gzip library with a simpler interface.
package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compress creates a compressed byte slice from a decompressed slice.
func Compress(src io.Reader) (io.Reader, int64, error) {
	dst := new(bytes.Buffer)

	writer := gzip.NewWriter(dst)
	defer writer.Close()

	len, err := io.Copy(writer, src)
	return dst, len, err
}

// Decompress creates a decompressed byte slice from a compressed slice.
// Forward errors from gzip's NewReader and io's Copy.
func Decompress(src io.Reader) (io.Reader, error) {
	dst := new(bytes.Buffer)

	reader, err := gzip.NewReader(src)
	if err == nil {
		_, err = io.Copy(dst, reader)
	}
	return dst, err
}
