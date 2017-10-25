package tar_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/skotchpine/xvm/util/tar"
)

func TestArchive(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvm-archive-test")
	path := filepath.Join(root, "pack")
	sub := filepath.Join(path, "subdirectory")
	norm := filepath.Join(path, "norm")

	if err := os.MkdirAll(sub, 0777); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(norm, []byte("key value"), 0777); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(root)

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

	if _, err := os.Stat(sub); err != nil {
		t.Errorf("The subdirectory %s was not restored form archive. Err: %s", sub, err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("The package directory %s was not restored from archive. Err: %s", path, err)
	}
}
