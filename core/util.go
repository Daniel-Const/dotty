package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type CopyError struct {
	message string
}

func (e *CopyError) Error() string {
	return e.message
}

/*
 * Copy 'from' file and place it at 'to' path.
 * Preserves file permissions
 */
func copyFile(src string, dest string) error {
	// Make destination dirs if not exists
	dirpath := filepath.Dir(dest)
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		return err
	}

	// Read dot file contents
	file, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Preserve file permissions
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Write dot file contents
	err = os.WriteFile(dest, file, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}

/*
 * Copy full directory tree from source to dest
 */
func copyDir(src string, dest string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		// log.Fatalf("Expecting directory at %s", src)
		return &CopyError{fmt.Sprintf("Expecting directory at %s", src)}
	}

	// Create path to the directory
	err = os.MkdirAll(dest, stat.Mode())
	if err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for i := range files {
		srcPath := filepath.Join(src, files[i].Name())
		destPath := filepath.Join(dest, files[i].Name())

		var err error
		if files[i].IsDir() {
			// Recursively copy next dir
			err = copyDir(srcPath, destPath)
		} else {
			err = copyFile(srcPath, destPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
