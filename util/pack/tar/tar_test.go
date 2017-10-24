package tar_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"."
)

func TestArchive(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvm-archive-test")
	path := filepath.Join(root, "pack")
	sub := filepath.Join(path, "subdirectory")
	dot := filepath.Join(path, ".hide")
	norm := filepath.Join(path, "norm")

	if err := os.MkdirAll(sub, 0777); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(dot, []byte("hidden"), 0777); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(norm, []byte("key value"), 0777); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(root)

	if _, err := tar.Archive(norm); err == nil {
		t.Error("Attempting to archive something other than a directory did not fail")
	}

	dist, err := tar.Archive(path)
	if err != nil {
		t.Error(err)
	}

	os.RemoveAll(path)

	if err := tar.Unarchive(root, dist); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(norm); err != nil {
		t.Errorf("The normal file %s was not restored from archive. Err: %s", norm, err)
	}

	if _, err := os.Stat(dot); err == nil {
		t.Errorf("The dotfile %s was restored from archive. Err: %s", dot, err)
	}

	if _, err := os.Stat(sub); err != nil {
		t.Errorf("The subdirectory %s was not restored form archive. Err: %s", sub, err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("The package directory %s was not restored from archive. Err: %s", path, err)
	}
}
