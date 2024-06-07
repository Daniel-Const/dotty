package core

import (
	"log"
	"os"
	"path/filepath"
)

// Represents a Dot file
type Dot struct {
    Path        string  // Current path of the file
    DeployPath  string  // Where we need to put this file (From the dotty.map)
    IsDir       bool
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
        dot.DeployPath = newDeployPath
    }

    // Determiner if dot is a directory
    if file, err := os.Stat(dot.Path); err != nil {
        // Try dest instead
        if file, err := os.Stat(dot.DeployPath); err != nil {
            // Both source and dest don't exist
            log.Printf("Src: %s, Dest: %s", dot.Path, dot.DeployPath)
            log.Fatalf("Source & Dest paths don't exist: %s ", err.Error())
        } else {
            dot.IsDir = file.Mode().IsDir()
        }
    } else {
        dot.IsDir = file.Mode().IsDir()
    }

    return &dot
}
