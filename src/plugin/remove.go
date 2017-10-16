package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"../utils"
)

func remove(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("Too few arguments")
	} else if len(args) > 4 {
		return fmt.Errorf("Too many arguments")
	}

	path := filepath.Join("plugins", args[3])
	absPath, err := utils.GetXvmSubDir(path)
	if err != nil {
		return err
	}

	return os.RemoveAll(absPath)
}
