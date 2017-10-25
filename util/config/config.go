// Package config implements a minimal key-value text store
// with parsing and a standard error.
package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// Config files are standardly formatted. ErrIllFormatted is used whenever
// the standard format is not used.
var (
	ErrIllFormatted = errors.New("Ill-formatted config")
)

// Read a key-value config string. Defer to Read function.
func ReadString(s string) (cfg map[string]string, err error) {
	return Read(bytes.NewBufferString(s))
}

// Read a key-value config buffer.
//
// Return ErrIllFormatted if the buffer does not conform to these two rules:
// 1) Keys are one word followed by whitespace and a value.
// 2) Values can contain any value but a newline.
//
// Return error from os.Open if not nil.
func Read(r io.Reader) (cfg map[string]string, err error) {
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
