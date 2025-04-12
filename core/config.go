package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type DottyConfig struct {
	Profiles []string
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
		// Failed to read path: Use a default config
		config := DottyConfig{}
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

	config := DottyConfig{}
	config.Profiles = paths
	return &config, nil
}
