package core

import (
	"log"
	"os"
	"path/filepath"
)

/*
 * Copy 'from' file and place it at 'to' path.
 * Preserves file permissions
 */
func copyFile(src string, dest string) error {
	// Make dirs if not exist
	dirpath := filepath.Dir(src)
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Read dot file contents
	file, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	// Preserve file permissions
	stat, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	// Write dot file contents
	err = os.WriteFile(dest, file, stat.Mode())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

/*
 * Copy full directory tree from source to dest
 */
func copyDir(src string, dest string) error {
	stat, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		log.Fatalf("Expecting directory at %s", src)
	}

	// Create path to the directory
	err = os.MkdirAll(dest, stat.Mode())
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
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
