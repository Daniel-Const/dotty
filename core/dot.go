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

func NewDot(path string, deployPath string, isDir bool) *Dot {
    dot := Dot{path, deployPath, isDir}

    // Expand ~ to full path
    if deployPath[0] == '~' {
        base, err := os.UserHomeDir()
        if err != nil {
            log.Fatal(err)
        }
        newDeployPath := filepath.Join(base, deployPath[2:])
        dot.DeployPath = newDeployPath
    }

    return &dot
}
