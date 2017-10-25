package util

import (
	"io"
	"os"
	"os/exec"

	"github.com/skotchpine/xvm/util/keyval"
)

// Standard file permissions and modes.
const (
	ModeClobber = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	PermPublic  = 0777
)

// Aggregate errors from listing a directory's entries.
func DirNames(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return dir.Readdirnames(0)
}

// Aggregate errors from reading a key-val map from file with keyval's Read.
func ReadMap(path string) (conf map[string]string, err error) {
	var file *os.File
	if file, err = os.Open(path); err == nil {
		conf, err = keyval.Parse(file)
	}
	return
}

// Aggregate errors from writing a key-val map to file with keyval's Write.
func WriteMap(path string, conf map[string]string) error {
	file, err := os.OpenFile(path, ModeClobber, PermPublic)
	if err != nil {
		return err
	}

	reader, err := keyval.NewReader(conf)
	if err == nil {
		_, err = io.Copy(file, reader)
	}
	return err
}

// Check if a file exists, discarding os's FileInfo.
func NotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// Execute a command, printing to stdout and stderr.
func Cmd(path string, arg ...string) error {
	c := exec.Command(path, arg...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
