package core

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DottyConfig struct {
	Profiles []string
	Path     string
}

func (c *DottyConfig) AddProfile(path string) {
	for _, p := range c.Profiles {
		if p == path {
			return
		}
	}
	c.Profiles = append(c.Profiles, path)
}

func (c *DottyConfig) Save() error {
	if c.Path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(c.Path), 0755); err != nil {
		return err
	}
	f, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, p := range c.Profiles {
		if _, err := f.WriteString(p + "\n"); err != nil {
			return err
		}
	}
	return nil
}

/* TODO: Extend config file
 * - Backups / Backup path
 * - Auto save new profiles to config
 * - Extend config file parsing for these options
 */

/*
 * Load a dotty configuration from a path
 * If no path is provided - use default values
 */
func LoadConfig(configPath string) (*DottyConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		config := DottyConfig{Path: configPath}
		return &config, err
	}

	paths := []string{}
	defer file.Close()

	// Scan map and create dots line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Skip if line is a comment
		checkIsComment := strings.Replace(line, " ", "", -1)
		if checkIsComment[0] == '#' {
			continue
		}

		p, err := processPath(line)
		if err != nil {
			log.Fatal(err)
		}
		paths = append(paths, p)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	config := DottyConfig{Profiles: paths, Path: configPath}
	return &config, nil
}
