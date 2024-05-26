package core

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)


type Profile struct {
    Name        string
    Os          string
    Location    string
    Dots        []*Dot
}

// Create a Profile struct from a name
func NewProfile(name string) *Profile {
    profile := Profile{}
    profile.Name = name
    profile.Location = "./profiles/" + name

    return &profile
}

func (p *Profile) Load() *Profile {
    var dots []*Dot

    // Read .map file and create a slice of dots
    file, err := os.Open(p.Location + "/dotty.map")
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
        toPath   := strings.Trim(parts[1], " ")
        fromPath := p.Location+"/"+fileName
        
        // Determine if directory or file
        file, err := os.Stat(fromPath)
        if err != nil {
            log.Fatal(err)
        }

        dots = append(dots, &Dot{fromPath, toPath, file.Mode().IsDir()})
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    p.Dots = dots

    return p
}

func (p *Profile) Deploy() error {
    for i := range p.Dots {
        dot := p.Dots[i]
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
func copyFile(from string, to string) error {
    
    // Make dirs if not exist
    dirpath := filepath.Dir(to)
    if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
        log.Fatal(err)
    }

    // Read dot file contents
    file, err := os.ReadFile(from)
    if err != nil {
        log.Fatal(err)
    }
    
    // Preserve file permissions
    stat, err := os.Stat(from)
    if err != nil {
        log.Fatal(err)
    }

    // Write dot file contents
    err = os.WriteFile(to, file, stat.Mode())
    if err != nil {
        log.Fatal(err)
    }

    return nil 
}

/*
 * Copy full directory tree at 'from' and place it at 'to'
 */
func copyDir(from string, to string) error {
    stat, err := os.Stat(from)
    if err != nil {
        log.Fatal(err)
    }

    if !stat.IsDir() {
        log.Fatalf("Expecting directory at %s", from)
    }
    
    // Create path to the directory
    err = os.MkdirAll(to, stat.Mode())
    if err != nil {
        log.Fatal(err)
    }

    files, err := os.ReadDir(from)
    for i := range files {
        fromPath := filepath.Join(from, files[i].Name())
        toPath := filepath.Join(to, files[i].Name())
        
        // Recursively copy next dir
        if files[i].IsDir() {
            err = copyDir(fromPath, toPath)
            if err != nil {
                return err
            }
        } else {
            // Copy the file
            err = copyFile(fromPath, toPath)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

