package core

import (
	"bufio"
	"errors"
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

func NewProfile(path string) *Profile {
    profile := Profile{}
    profile.Name = filepath.Base(path)
    profile.Location = path 
    return &profile
}

func (p *Profile) LoadMap() (*Profile, error) {
    var dots []*Dot

    // Read .map file and create a slice of dots
    file, err := os.Open(filepath.Join(p.Location, "dotty.map"))
    if err != nil {
        return p, err
    }
    defer file.Close()
    
    // Scan map and create dots line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue
        }
        parts := strings.Split(line, ":")
        if len(parts) < 2 {
           log.Println(line)
            return p, errors.New("Err: Map file incorrectly formatted")
        }

        fileName := strings.Trim(parts[0], " ")
        destPath   := strings.Trim(parts[1], " ")
        sourcePath := filepath.Join(p.Location, fileName)
        dots = append(dots, NewDot(sourcePath, destPath))
    }

    if err := scanner.Err(); err != nil {
        return p, err
    }

    p.Dots = dots
    return p, nil
}

/*
 * Copy files at destination paths into a profile
 */
func (p *Profile) Load() error {
    for i := range p.Dots {
        err := p.Dots[i].Load()
        if err != nil {
            return err
        }
    }
    return nil
}

/*
 * Copy all of the dotfiles to the locations in the map file
 */
func (p *Profile) Deploy() error {
    for i := range p.Dots {
        if err := p.Dots[i].Backup(); err != nil {
            return err 
        }

        if err := p.Dots[i].Deploy(); err != nil {
            return err 
        }
    }

    return nil
}


