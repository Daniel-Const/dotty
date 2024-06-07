package core

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)


type Profile struct {
    Name        string
    Os          string
    Location    string
    Dots        []*Dot
}

func NewProfile(path string) *Profile {
    profile := Profile{}
    profile.Name = filepath.Base(path)
    profile.Location = path 
    return &profile
}

func (p *Profile) LoadMap() *Profile {
    var dots []*Dot

    // Read .map file and create a slice of dots
    file, err := os.Open(filepath.Join(p.Location, "dotty.map"))
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    // Scan map and create dots line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, ":")
        if len(parts) < 2 {
            log.Fatal("Err: Map file incorrectly formatted")
        }

        fileName := strings.Trim(parts[0], " ")
        destPath   := strings.Trim(parts[1], " ")
        sourcePath := filepath.Join(p.Location, fileName)
        dots = append(dots, NewDot(sourcePath, destPath))
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    p.Dots = dots

    return p
}

/*
 * Copy files at destination paths into a profile
 */
func (p *Profile) Load() error {
    for i := range p.Dots {
        dot := p.Dots[i]
        log.Printf("Loading dotfile into profile: %s", dot.DeployPath)
        var err error
        if dot.IsDir {
            err = copyDir(dot.DeployPath, dot.Path)
        } else {
            err = copyFile(dot.DeployPath, dot.Path)
        }
        if err != nil {
            log.Fatal(err)
        }
    }
    return nil
}

/*
 * Copy all of the dotfiles to the locations in the map file
 */
func (p *Profile) Deploy() error {
    backupName := time.Now().Format("2006-01-02_15:04:05.000000")
    backupRoot := filepath.Join("./backup", backupName)
    for i := range p.Dots {
        dot := p.Dots[i]
        
        // Create backups if file / dir exists already
        if stat, err := os.Stat(dot.DeployPath); err == nil {
            backupPath := filepath.Join(backupRoot, filepath.Base(dot.DeployPath))
            log.Printf("Creating backup for %s", dot.DeployPath)

            if err := os.MkdirAll(backupRoot, os.ModePerm); err != nil {
                log.Fatal(err)
            }

            if stat.IsDir() {
                err := copyDir(dot.DeployPath, backupPath)
                if err != nil {
                    log.Fatal(err)
                }
            } else {
                err := copyFile(dot.DeployPath, backupPath)
                if err != nil {
                    log.Fatal(err)
                }
            }
        }
        
        // Copy files to their destination
        if dot.IsDir {
            if err := copyDir(dot.Path, dot.DeployPath); err != nil {
                log.Fatal(err)
            }
        } else {
            if err := copyFile(dot.Path, dot.DeployPath); err != nil {
                log.Fatal(err)
            }
        }

        log.Printf("Copied %s to %s", dot.Path, dot.DeployPath)
    }

    return nil
}

/*
 * Copy 'from' file and place it at 'to' path.
 * Preserves file permissions 
 */
func copyFile(source string, dest string) error {
    // Make dirs if not exist
    dirpath := filepath.Dir(source)
    if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
        log.Fatal(err)
    }

    // Read dot file contents
    file, err := os.ReadFile(source)
    if err != nil {
        log.Fatal(err)
    }
    
    // Preserve file permissions
    stat, err := os.Stat(source)
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
func copyDir(source string, dest string) error {
    stat, err := os.Stat(source)
    if err != nil {
        log.Fatal(err)
    }

    if !stat.IsDir() {
        log.Fatalf("Expecting directory at %s", source)
    }
    
    // Create path to the directory
    err = os.MkdirAll(dest, stat.Mode())
    if err != nil {
        log.Fatal(err)
    }

    files, err := os.ReadDir(source)
    for i := range files {
        sourcePath := filepath.Join(source, files[i].Name())
        destPath := filepath.Join(dest, files[i].Name())
        
        // Recursively copy next dir
        if files[i].IsDir() {
            err = copyDir(sourcePath, destPath)
            if err != nil {
                return err
            }
        } else {
            // Copy the file
            err = copyFile(sourcePath, destPath)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
