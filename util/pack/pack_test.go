package pack_test

import (
	"os"
	"testing"

	"github.com/skotchpine/xvm/util/pack"
)

func TestPackContext(t *testing.T) {
	config := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}

	env := map[string]string{
		"XVM_PULL_PATH":    "qwer",
		"XVM_PULL_CONFIG":  "",
		"XVM_PULL_VERSION": "zxcv",
	}

	for key, val := range config {
		env["XVM_PULL_CONFIG"] = env["XVM_PULL_CONFIG"] + key + " " + val + "\n"
	}

	for key, val := range env {
		if err := os.Setenv(key, val); err != nil {
			t.Error(err)
		}
	}

	ctx, err := pack.Context()
	if err != nil {
		t.Error(err)
	}

	if ctx.Path != env["XVM_PULL_PATH"] {
		t.Error("Failed to get XVM_PULL_PATH from the environment")
	}
	if ctx.Version != env["XVM_PULL_VERSION"] {
		t.Error("Failed to get XVM_PULL_VERSION from the environment")
	}

	for key, expected := range config {
		if actual, ok := ctx.Config[key]; !ok {
			t.Errorf("Missing %s from config", key)
		} else if actual != expected {
			t.Errorf("Expected %s to have value %s, but got %s", key, expected, actual)
		}
	}
}
