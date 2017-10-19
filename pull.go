package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Gets a new dist.zip, unzips it and truncates conflicting files.
func main() {
	destDir := os.Getenv("XVM_PULL_DESTDIR") // where X should be installed
	version := os.Getenv("XVM_PULL_VERSION") // the name of the version requested
	content := os.Getenv("XVM_PULL_CONTENT") // the content of the version's file

	fail := func(msg string, etc ...interface{}) {
		fmt.Fprintf(os.Stderr, msg+"\n", etc...)
		fmt.Printf("Failed to install plugin: %s\n", version)
		os.Exit(1)
	}

	failIf := func(err error) {
		if err != nil {
			fail(err.Error())
		}
	}

	// Content is the url of a dist.zip file.
	resp, err := http.Get(content)
	failIf(err)
	if resp.StatusCode != 200 {
		fail("Failed to get dist.zip from %s", content)
	}

	// Writing dist.zip to disk, because zip.OpenReader takes a path.
	body, err := ioutil.ReadAll(resp.Body)
	failIf(err)

	distPath := filepath.Join(destDir, "dist.zip")
	failIf(ioutil.WriteFile(distPath, body, 0755))
	defer os.Remove(distPath)

	// Copying everything in the archive to disk.
	distReader, err := zip.OpenReader(distPath)
	failIf(err)

	for _, zipfile := range distReader.File {
		// The file's destination.
		realPath := filepath.Join(destDir, zipfile.Name)

		// Copying directories will be non-destructive.
		if zipfile.FileInfo().IsDir() {
			os.MkdirAll(realPath, zipfile.Mode())
			continue
		}

		fileReader, err := zipfile.Open()
		failIf(err)
		defer fileReader.Close()

		// Copying files will be destructive.
		mask := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		fileWriter, err := os.OpenFile(realPath, mask, zipfile.Mode())
		failIf(err)
		defer fileWriter.Close()

		_, err = io.Copy(fileWriter, fileReader)
		failIf(err)
	}

	fmt.Printf("Successfully installed plugin: %s\n", version)
}
