package keyval_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/skotchpine/xvm/util/keyval"
)

var (
	expecteds = map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}
	lineDelims   = []string{"\n", "\r\n"}
	keyValDelims = []string{"\t", " ", "    "}

	inputs []io.Reader
)

func init() {
	for _, keyValDelim := range keyValDelims {
		var input []byte
		for _, lineDelim := range lineDelims {
			for key, val := range expecteds {
				input = append(input, (key + keyValDelim + val + lineDelim)...)
			}
		}
		inputs = append(inputs, bytes.NewBuffer(input))
	}
}

func TestNewReader(t *testing.T) {
	reader, err := keyval.NewReader(expecteds)
	if err != nil {
		t.Error(err)
	}

	actuals, err := keyval.Parse(reader)
	if err != nil {
		t.Error(err)
	}

	for key, expected := range expecteds {
		actual, ok := actuals[key]
		if !ok {
			t.Errorf("Key '%s' not parsed", key)
		}
		if expected != actual {
			t.Errorf(
				"Expected key '%s' to be '%s', but got '%s'",
				key, expected, actual,
			)
		}
	}
}

func TestParse(t *testing.T) {
	if _, err := keyval.ParseString("hi"); err != keyval.ErrIllFormatted {
		t.Error("An empty key did not cause ErrIllFormatted")
	}

	if _, err := keyval.ParseString("hi"); err != keyval.ErrIllFormatted {
		t.Error("A key without a value did not cause ErrIllFormatted")
	}

	for _, input := range inputs {
		actuals, err := keyval.Parse(input)
		if err != nil {
			t.Error(err)
		}

		for key, expected := range expecteds {
			actual, ok := actuals[key]
			if !ok {
				t.Errorf("Key '%s' not parsed", key)
			}
			if expected != actual {
				t.Errorf(
					"Expected key '%s' to be '%s', but got '%s'",
					key, expected, actual,
				)
			}
		}
	}
}
