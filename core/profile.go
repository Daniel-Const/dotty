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

        isDir := file.Mode().IsDir()
        dot := Dot{fromPath, toPath, isDir}
        dots = append(dots, &dot)
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

        // Make dirs if not exist
        dirpath := filepath.Dir(dot.DeployPath)
        if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
            log.Fatal(err)
        }

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

func copyFile(from string, to string) error {
    log.Println("Copying file")

    // Read dot file contents
    file, err := os.ReadFile(from)
    if err != nil {
        log.Fatal(err)
    }

    // Write dot file contents
    err = os.WriteFile(to, file, 0644)
    if err != nil {
        log.Fatal(err)
    }

    return nil 
}

func copyDir(from string, to string) error {
    // TODO: Implement
    return nil
}

