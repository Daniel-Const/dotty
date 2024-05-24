package main

import (
    "log"
    "strings"
    "bufio"
    "os"
)


type Profile struct {
    Name        string
    Os          string
    Location    string
    dots        []*Dot
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

        dot := Dot{p.Location+"/"+fileName, toPath}
        dots = append(dots, &dot)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    p.dots = dots

    return p
}



