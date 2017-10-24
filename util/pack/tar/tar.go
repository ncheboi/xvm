package tar

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

func Archive(root string) ([]byte, error) {
	buf := new(bytes.Buffer)
	archive := tar.NewWriter(buf)

	if info, err := os.Stat(root); err != nil {
		return buf.Bytes(), err
	} else if !info.IsDir() {
		return buf.Bytes(), errors.New("Attempting to archive something other than a directory")
	}

	var err error
	if err = archiveFile(archive, root, filepath.Base(root)); err == nil {
		err = archive.Close()
	}
	return buf.Bytes(), err
}

func Unarchive(dstPath string, dist []byte) error {
	archive := tar.NewReader(bytes.NewBuffer(dist))
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if info.IsDir() {
			dst := filepath.Join(dstPath, header.Name)
			if err = os.MkdirAll(dst, info.Mode()); err != nil {
				return err
			}
			continue
		}

		mask := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		dst := filepath.Join(dstPath, header.Name)
		file, err := os.OpenFile(dst, mask, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(file, archive); err != nil {
			return err
		}
	}
	return nil
}

func archiveDir(archive *tar.Writer, abs, rel string) error {
	file, err := os.Open(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		absEntry := filepath.Join(abs, entry)
		relEntry := path.Join(rel, entry)
		if err := archiveFile(archive, absEntry, relEntry); err != nil {
			return err
		}
	}
	return nil
}

func archiveFile(archive *tar.Writer, abs, rel string) error {
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if info.Name()[0] == '.' {
		return nil
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = rel
	if err = archive.WriteHeader(header); err != nil {
		return err
	}

	if info.IsDir() {
		return archiveDir(archive, abs, rel)
	}

	file, err := os.Open(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(archive, file)
	return err
}
