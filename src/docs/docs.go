package docs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"../utils"
)

// Print a file in $XVMPATH/doc to Stdout.
func Print(name string) error {
	dir, err := utils.GetXvmSubDir("doc")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Print(string(buf))

	return nil
}

