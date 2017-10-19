package plugin_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"."
)

var (
	tmp = os.TempDir()

	// The distFiles which are also directories.
	dirs = []string{
		filepath.Join(tmp, "plugin", "windows"),
		filepath.Join(tmp, "plugin", "unix"),
		filepath.Join(tmp, "plugin", "src"),
	}

	// Some files that should be covered by distFiles.
	files = map[string]string{
		filepath.Join(tmp, "plugin", "windows", "plugin.bat"): "echo %*\n",
		filepath.Join(tmp, "plugin", "unix", "plugin"): "#!/bin/sh\n\necho $@\n",
		filepath.Join(tmp, "plugin", "src", "s1.c"): "int main() \n",
		filepath.Join(tmp, "plugin", "available"): "1\n2\n3\n",
	}

	// Some directories that should *not* be affected by Plugin.Update.
	userdirs = []string{
		filepath.Join(tmp, "plugin", ".git"),
		filepath.Join(tmp, "plugin", "installed", "v1"),
	}

	// Some files that should *not* be affected by Plugin.Update.
	userfiles = map[string]string{
		filepath.Join(tmp, "plugin", ".git", "config"): "[user]\n\tname = skotchpine",
	}
)

func TestDist(t *testing.T) {
	// Make directories for all testing conditions.
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Errorf("Failed to create plugin directory: %v", err)
		}
	}

	for _, dir := range userdirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Errorf("Failed to create user directory: %v", err)
		}
	}

	for file, contents := range files {
		if err := ioutil.WriteFile(file, []byte(contents), 0755); err != nil {
			t.Errorf("Failed to create plugin file: %v", err)
		}
	}

	for file, contents := range userfiles {
		if err := ioutil.WriteFile(file, []byte(contents), 0755); err != nil {
			t.Errorf("Failed to create plugin file: %v", err)
		}
	}

	p := plugin.Plugin{
		Path: filepath.Join(tmp, "plugin"),
	}
	defer os.RemoveAll(p.Path)

	t.Run("plugin.Bundle should create a new dist.zip", func(t *testing.T) {
		zip, err := p.Bundle()
		if err != nil {
			t.Error(err)
		}
		if _, err  = os.Stat(zip); err != nil {
			t.Error(err)
		}
	})

	// Move dist.zip to a new place and serve it over http
	// where the plugin expects it.
	distPath := filepath.Join(p.Path, "dist.zip")
	servPath := filepath.Join(p.Path, "serv.zip")
	if err := os.Rename(distPath, servPath); err != nil {
		t.Errorf("Failed to move dist file for http server: %v", err)
	}
	defer os.RemoveAll(servPath)

	handler := http.FileServer(http.Dir(p.Path))
	serv := httptest.NewServer(handler)
	defer serv.Close()

	p.Url = serv.URL + "/serv.zip"

//	for file, _ := range files {
//		if err := os.Remove(file); err != nil {
//			t.Errorf("Failed to remove plugin file: %v", err)
//		}
//	}
//
//	for _, dir := range dirs {
//		if err := os.RemoveAll(dir); err != nil {
//			t.Errorf("Failed to remove plugin directory: %v", err)
//		}
//	}

	t.Run("plugin.Update should install a new dist.zip", func(t *testing.T) {
		if err := p.Update(); err != nil {
			t.Error(err)
		}

		for _, dir := range dirs {
			info, err := os.Stat(dir)
			if err != nil || !info.IsDir() {
				t.Errorf("Failed to transfer plugin directory: %s", dir)
			}
		}

		for file, expected := range files {
			actual, err := ioutil.ReadFile(file)
			if err != nil {
				if os.IsNotExist(err) {
					t.Errorf("Failed to transfer plugin file: %s", file)
				} else {
					t.Errorf("Failed to read plugin file: %s", file)
				}
			} else if string(actual) != expected {
				t.Errorf("Failed to transfer plugin file: %s", file)
			}
		}

		for _, dir := range userdirs {
			info, err := os.Stat(dir)
			if err != nil || !info.IsDir() {
				t.Errorf("Failed to avoid user directory: %s", dir)
			}
		}

		for file, expected := range userfiles {
			actual, err := ioutil.ReadFile(file)
			if err != nil {
				if os.IsNotExist(err) {
					t.Errorf("Failed to avoid user file: %s", file)
				} else {
					t.Errorf("Failed to read user file: %s", file)
				}
			} else if string(actual) != expected {
				t.Errorf("Failed to avoid user file: %s", file)
			}
		}
	})
}
