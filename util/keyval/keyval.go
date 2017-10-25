// Package config implements a minimal key-value text store
// with parsing and a standard error.
package keyval

import (
	"bufio"
	"bytes"
	"io"
)

// All files will be ended with \n and all key-val pairs will be separated by
// one space by default.
const (
	LineDelim   = "\n"
	KeyValDelim = " "
)

// Implement io's Reader from a config. Forward errors from io operations.
func NewReader(cfg map[string]string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	for key, val := range cfg {
		if err := writeLine(key, val, buf); err != nil {
			return buf, err
		}
	}
	return buf, nil
}

// Write one key-val pair to a buffer.
func writeLine(key, val string, buf io.Writer) (err error) {
	var line []byte
	if val == "" {
		line = []byte(key + LineDelim)
	} else {
		line = []byte(key + KeyValDelim + val + LineDelim)
	}
	_, err = buf.Write(line)
	return
}

// Read a key-value config string. Delegate to Read function.
func ParseString(s string) (cfg map[string]string, err error) {
	return Parse(bytes.NewBufferString(s))
}

// Parse a key-value config buffer. Forward errors from os.Open if not nil.
func Parse(r io.Reader) (cfg map[string]string, err error) {
	// Use bufio's Scanner and ScanLines to split the file into lines.
	// A side-effect of this is that all \r characters will be stripped,
	// so any \r character must be accompanied by a \n to end a line.
	cfg = make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		key, val := parseLine(scanner.Bytes())
		cfg[key] = val
	}
	return
}

func parseLine(buf []byte) (string, string) {
	n := len(buf) - 1
	lead := 0
	for lead < n && buf[lead] == ' ' || buf[lead] == '\t' {
		lead++
	}
	cursor := lead

	for cursor < n && buf[cursor] != ' ' && buf[cursor] != '\t' {
		cursor++
	}
	if cursor == n {
		return string(buf[lead : cursor+1]), ""
	}
	key := string(buf[lead:cursor])

	for cursor < n && buf[cursor] == ' ' || buf[cursor] == '\t' {
		cursor++
	}
	return key, string(buf[cursor:])
}
