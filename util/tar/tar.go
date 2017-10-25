// Package tar implements golang's tar library with a simpler interface.
package tar

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Archive creates an archived byte slice from the path of a directory.
// Forward errors from operating system queries, input/output and archive operatinos.
func Archive(root string) (io.Reader, error) {
	dst := new(bytes.Buffer)
	archive := tar.NewWriter(dst)

	// Write files to the archive starting with the root directory. If the file
	// to be written is a directory, write all of its entries to the archive.
	var err error
	if err = archiveWriteFile(archive, root, filepath.Base(root)); err == nil {
		err = archive.Close()
	}
	return dst, err
}

// Write to archive the file at the absolute path abs to the relative path rel.
func archiveWriteFile(archive *tar.Writer, abs, rel string) error {
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}

	// Write a new header to the archive with the path rel.
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = rel
	if err = archive.WriteHeader(header); err != nil {
		return err
	}

	// Write directories recursively and return; no content to copy.
	if info.IsDir() {
		return archiveWriteEntries(archive, abs, rel)
	}

	// Copy content from all non-directories to archive.
	file, err := os.Open(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(archive, file)
	return err
}

// For each entry in the directory at absolute path abs, write a file to archive
// with the relative prefix rel.
func archiveWriteEntries(archive *tar.Writer, abs, rel string) error {
	file, err := os.Open(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	// filepath's Join is used for absolute paths with the operating system's
	// path separator (/ or \), but path's Join is used for relative tar paths,
	// because the unix path separator (/) is used in tar archives.
	for _, entry := range entries {
		absEntry := filepath.Join(abs, entry)
		relEntry := path.Join(rel, entry)
		if err := archiveWriteFile(archive, absEntry, relEntry); err != nil {
			return err
		}
	}
	return nil
}

// Unarchive creates a new directory inside abs with the contents of archive,
// because Archive writes files relative to and including a root directory.
// Any paths which conflict may be truncated or have their permission bits reset.
//
// Forward errors from operating system queries, input/output and archive operations.
func Unarchive(abs string, src io.Reader) error {
	archive := tar.NewReader(src)

	// Get file info for each header in the archive until EOF.
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		info := header.FileInfo()

		// Skip after creation if the file is a directory.
		if info.IsDir() {
			a := filepath.Join(abs, header.Name)
			if err = os.MkdirAll(a, info.Mode()); err != nil {
				return err
			}
			continue
		}

		// Get the file handle for all non-directories, and overwrite content.
		mask := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		a := filepath.Join(abs, header.Name)
		file, err := os.OpenFile(a, mask, info.Mode())
		if err == nil {
			_, err = io.Copy(file, archive)
		}
		file.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
