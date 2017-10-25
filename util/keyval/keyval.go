// Package config implements a minimal key-value text store
// with parsing and a standard error.
package keyval

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// Config files are standardly formatted. ErrIllFormatted is used whenever
// the standard format is not used.
var (
	ErrIllFormatted = errors.New("Ill-formatted key-value file")
)

// Implement io's Reader from a config. Forward errors from io operations.
func NewReader(cfg map[string]string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	var err error
	for key, val := range cfg {
		if _, err = buf.Write([]byte(key + " " + val + "\n")); err != nil {
			break
		}
	}
	return buf, err
}

// Read a key-value config string. Delegate to Read function.
func ParseString(s string) (cfg map[string]string, err error) {
	return Parse(bytes.NewBufferString(s))
}

// Parse a key-value config buffer.
//
// Return ErrIllFormatted if the buffer does not conform to these two rules:
// 1) One key per line, starting at the first character.
// 2) Everything after the first whitespace and before the next newline is the value.
//
// Return error from os.Open if not nil.
func Parse(r io.Reader) (cfg map[string]string, err error) {
	// Use bufio's Scanner and ScanLines to split the file into lines.
	// A side-effect of this is that all \r characters will be stripped,
	// so any \r character must be accompanied by a \n to end a line.
	cfg = make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		// If the first character is a whitespace, the key is empty.
		if line[0] == ' ' || line[0] == '\t' {
			err = ErrIllFormatted
			return
		}

		// Identify the limits of the first whitespace string
		// anywhere but the first character.
		start, stop := 0, 0
		for i, b := range line {
			if b == ' ' || b == '\t' {
				if i != 0 && start == 0 {
					start = i
				}
			} else {
				if start != 0 {
					stop = i
					break
				}
			}
		}

		// Slice the key and value out of line, using the whitespace as the delimiter.
		// If either bound of the whitespace was not found, the file is Ill-formatted.
		if start != 0 && stop != 0 {
			cfg[string(line[:start])] = string(line[stop:])
		} else {
			err = ErrIllFormatted
		}
	}
	return
}
