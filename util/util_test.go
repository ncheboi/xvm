package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/skotchpine/xvm/util"
)

func TestDirNames(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvm-dirnames-test")

	dir := filepath.Join(root, "dir")
	if err := os.MkdirAll(dir, util.PermPublic); err != nil {
		t.Error(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}()

	names := []string{"a", "b", "c"}
	for _, name := range names {
		ioutil.WriteFile(filepath.Join(dir, name), []byte{}, util.PermPublic)
	}

	n, err := util.DirNames(dir)
	if err != nil {
		t.Error(err)
	}
	for _, expected := range names {
		missing := true

		for _, actual := range n {
			if actual == expected {
				missing = false
				break
			}
		}

		if missing {
			t.Errorf("Did not list file %s", expected)
		}
	}
}

func TestMap(t *testing.T) {
	path := filepath.Join(os.TempDir(), "xvm-test-map")
	expected := map[string]string{
		"key1": "val2",
		"key2": "val2",
	}

	err := util.WriteMap(path, expected)
	if err != nil {
		t.Error(err)
	}

	actual, err := util.ReadMap(path)
	if err != nil {
		t.Error(err)
	}

	for key, e := range expected {
		if a, ok := actual[key]; !ok {
			t.Errorf("%s not found after write and read", key)
		} else if a != e {
			t.Errorf("Expected %s to have value %s, but got %s", key, e, a)
		}
	}
}

func TestNotExist(t *testing.T) {
	path := filepath.Join(os.TempDir(), "xvm-test-dir-not-exist")

	if err := os.MkdirAll(path, util.PermPublic); err != nil {
		t.Error(err)
	} else if util.NotExist(path) {
		t.Error("NotExist returned true, but directory exists.")
	}

	if err := os.RemoveAll(path); err != nil {
		t.Error(err)
	} else if !util.NotExist(path) {
		t.Error("NotExist returned false, but directory exists.")
	}
}

func TestCmd(t *testing.T) {
	t.Skip()
}
