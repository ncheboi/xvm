package keyval_test

import (
	"bytes"
	"testing"

	"github.com/skotchpine/xvm/util/keyval"
)

var (
	lineDelims   = []string{"\n", "\r\n"}
	keyValDelims = []string{"\t", " ", "    "}
)

func notErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func compare(t *testing.T, expected, actual map[string]string) {
	for key, e := range expected {
		a, ok := actual[key]
		if !ok {
			t.Errorf("Key '%s' not parsed", key)
		}
		if a != e {
			t.Errorf("Expected value of '%s' to be '%s', but got '%s'", key, e, a)
		}
	}
}

func TestNewReader(t *testing.T) {
	expected := map[string]string{"key1": "val1", "key2": ""}

	reader, err := keyval.NewReader(expected)
	notErr(t, err)

	actual, err := keyval.Parse(reader)
	notErr(t, err)

	compare(t, expected, actual)
}

func TestLineDelims(t *testing.T) {
	expected := map[string]string{"key1": "val1", "key2": "val2"}

	for _, delim := range lineDelims {
		actual, err := keyval.ParseString("key1 val1" + delim + "key2 val2" + delim)
		notErr(t, err)

		compare(t, expected, actual)
	}
}

func TestKeyValDelims(t *testing.T) {
	expected := map[string]string{"key1": "val1"}

	for _, delim := range keyValDelims {
		buf := new(bytes.Buffer)
		buf.Write([]byte("key1" + delim + "val1\n"))
		buf.Write([]byte("key2" + delim + "val2\n"))

		actual, err := keyval.Parse(buf)
		notErr(t, err)

		compare(t, expected, actual)
	}
}

func TestLeadingWhitespace(t *testing.T) {
	expected := map[string]string{"key1": "val1"}

	buf := new(bytes.Buffer)
	buf.Write([]byte(" key1 val1\n"))

	actual, err := keyval.Parse(buf)
	notErr(t, err)

	compare(t, expected, actual)
}
