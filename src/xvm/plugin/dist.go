package plugin

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Bundle will only include these files.
var distFiles = []string{
	"available",
	"windows",
	"unix",
	"src",
}

// Creates a zip archive of all distribution files. Returns the archive path.
func (p *Plugin) Bundle() (string, error) {
	distPath := filepath.Join(p.Path, "dist.zip")
	distFile, err := os.Create(distPath)
	if err != nil {
		return "", err
	}
	defer distFile.Close()

	dist := zip.NewWriter(distFile)
	defer dist.Close()

	// Use bundleFile to write all files in distFiles to dist.zip recursively.
	for _, zipPath := range distFiles {
		realPath := filepath.Join(p.Path, zipPath)
		if err := bundleFile(dist, realPath, zipPath); err != nil {
			return "", err
		}
	}

	return distPath, nil
}

// Writes a file to a zip archive, recursing if the file is a non-empty directory.
func bundleFile(dist *zip.Writer, realPath, zipPath string) error {
	info, err := os.Stat(realPath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	if info.IsDir() {
		// Directories must end in a slash (per the zip spec).
		filepath.ToSlash(zipPath)
	} else {
		header.Method = zip.Deflate
		header.Name = zipPath
	}

	writer, err := dist.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(realPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file and return unless it's a directory.
	if !info.IsDir() {
		_, err = io.Copy(writer, file)
		return err
	}

	// Recurse for each entry in the directory.
	entries, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryRealPath := filepath.Join(realPath, entry)
		entryZipPath := filepath.Join(zipPath, entry)

		if err := bundleFile(dist, entryRealPath, entryZipPath); err != nil {
			return err
		}
	}

	return nil
}

// Get a new dist.zip and uncompress it, truncating existing distFiles.
func (p *Plugin) Update() error {
	resp, err := http.Get(p.Url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to get %s", p.Url)
	}

	// Writes dist.zip to disk, because zip.OpenReader takes a path.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	distPath := filepath.Join(p.Path, "dist.zip")
	if err := ioutil.WriteFile(distPath, body, 0755); err != nil {
		return err
	}
	defer os.Remove(distPath)

	reader, err := zip.OpenReader(distPath)
	if err != nil {
		return err
	}

	// Copy everything in the archive to disk.
	for _, zipfile := range reader.File {
		realPath := filepath.Join(p.Path, zipfile.Name)
		if zipfile.FileInfo().IsDir() {
			os.MkdirAll(realPath, zipfile.Mode())
			continue
		}

		reader, err := zipfile.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		mask := os.O_WRONLY|os.O_CREATE|os.O_TRUNC
		file, err := os.OpenFile(realPath, mask, zipfile.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(file, reader); err != nil {
			return err
		}
	}

	return nil
}
