package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"."
)

var (
	expecteds = map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}
	lineDelims   = []string{"\n", "\r\n"}
	keyValDelims = []string{"\t", " ", "    "}

	inputs [][]byte
)

func init() {
	for _, keyValDelim := range keyValDelims {
		var input []byte
		for _, lineDelim := range lineDelims {
			for key, val := range expecteds {
				input = append(input, (key + keyValDelim + val + lineDelim)...)
			}
		}
		inputs = append(inputs, input)
	}
}

func TestParse(t *testing.T) {
	configPath := "xvm-test-config"

	os.RemoveAll(configPath)
	defer os.RemoveAll(configPath)

	if _, err := config.Parse(configPath); err == nil {
		t.Error("The path of a nonexistent file did not cause error")
	}

	ioutil.WriteFile(configPath, []byte("hi"), 0777)
	if _, err := config.Parse(configPath); err != config.ErrIllFormatted {
		t.Error("A key without a value did not cause ErrIllFormatted")
	}

	ioutil.WriteFile(configPath, []byte(" hi"), 0777)
	if _, err := config.Parse(configPath); err != config.ErrIllFormatted {
		t.Error("An empty key did not cause ErrIllFormatted")
	}

	for _, input := range inputs {
		ioutil.WriteFile(configPath, input, 0777)

		actuals, err := config.Parse(configPath)
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
