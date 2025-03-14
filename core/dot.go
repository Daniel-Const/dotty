package core

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Represents a Dot file
type Dot struct {
	SrcPath  string // Current path of the file
	DestPath string // Where we need to put this file (From the dotty.map)
	IsDir    bool
}

func NewDot(path string, deployPath string) *Dot {
	dot := Dot{path, deployPath, false}

	// Expand ~ to full path
	if deployPath[0] == '~' {
		base, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		newDeployPath := filepath.Join(base, deployPath[2:])
		dot.DestPath = newDeployPath
	}

	// Determine if dot is a directory
	if file, err := os.Stat(dot.SrcPath); err != nil {
		// Try dest instead
		if file, err := os.Stat(dot.DestPath); err != nil {
			// Both source and dest don't exist
			log.Printf("Src: %s, Dest: %s", dot.SrcPath, dot.DestPath)
			log.Fatalf("Source & Dest paths don't exist: %s ", err.Error())
		} else {
			dot.IsDir = file.Mode().IsDir()
		}
	} else {
		dot.IsDir = file.Mode().IsDir()
	}

	return &dot
}

// Copy SrcPath files to DestPath
func (d *Dot) Deploy() error {
	var err error
	if d.IsDir {
		err = copyDir(d.SrcPath, d.DestPath)
	} else {
		err = copyFile(d.SrcPath, d.DestPath)
	}

	if err != nil {
		return err
	}

	log.Printf("Copied %s to %s", d.SrcPath, d.DestPath)
	return nil
}

// Copy DestPath files to SrcPath
func (d *Dot) Load() error {
	var err error
	if d.IsDir {
		err = copyDir(d.DestPath, d.SrcPath)
	} else {
		err = copyFile(d.DestPath, d.SrcPath)
	}

	if err != nil {
		return err
	}

	log.Printf("Loading dotfile into profile: %s", d.DestPath)
	return nil
}

// Backup DestPath files if they exist already
func (d *Dot) Backup() error {
	// TODO: Get backup path dir from viper config?
	backupName := time.Now().Format("2006-01-02_15:04:05.00")
	backupRoot := filepath.Join("backup", backupName)
	if stat, err := os.Stat(d.DestPath); err == nil {
		backupPath := filepath.Join(backupRoot, filepath.Base(d.DestPath))
		log.Printf("Creating backup for %s", d.DestPath)
		if err := os.MkdirAll(backupRoot, os.ModePerm); err != nil {
			return err
		}

		var err error
		if stat.IsDir() {
			err = copyDir(d.DestPath, backupPath)
		} else {
			err = copyFile(d.DestPath, backupPath)
		}

		if err != nil {
			return err
		}
	}

	log.Printf("Created backup for %s", d.DestPath)
	return nil
}
