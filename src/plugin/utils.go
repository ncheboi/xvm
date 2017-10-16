package plugin

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../utils"
)

// Gets a map of local plugins and their absolute paths.
func getLocalPlugins() (map[string]string, error) {
	dir, err := utils.GetXvmSubDir("plugins")
	if err != nil {
		return nil, err
	}

	plugins := make(map[string]string)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == ".gitkeep" {
			continue
		}

		path := filepath.Join(dir, name)

		if !entry.IsDir() {
			return nil, fmt.Errorf("The plugin %s at %s is not a directory", name, path)
		}

		plugins[name] = path
	}

	return plugins, nil
}

// Gets a map of remote plugins and their urls.
func getRemotePlugins() (map[string]string, error) {
	dir, err := utils.GetXvmSubDir("doc")
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "plugins")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	plugins := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("The plugins file %s is ill-formatted", path)
		}

		plugins[tokens[0]] = tokens[1]
	}

	return plugins, nil
}

