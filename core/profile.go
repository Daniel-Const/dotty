package core

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Profile struct {
	Name     string
	Os       string
	Location string
	Dots     []*Dot
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

		// Skip if line is a comment
		checkIsComment := strings.Replace(line, " ", "", -1)
		if checkIsComment[0] == '#' {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			log.Println(line)
			return p, errors.New("err: Map file incorrectly formatted")
		}

		fileName := strings.Trim(parts[0], " ")
		destPath := strings.Trim(parts[1], " ")
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
			fmt.Println(err)
		} else {
			fmt.Println("Loaded ", p.Dots[i].DestPath, " Into ", p.Dots[i].SrcPath)
		}
	}
	return nil
}

/*
 * Copy all of the dotfiles to the locations in the map file
 */
func (p *Profile) Deploy() error {
	for i := range p.Dots {
		if err := p.Dots[i].Deploy(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Deployed ", p.Name, " To ", p.Location)
		}
	}

	return nil
}

func (p *Profile) Print() {
	// var style = lipgloss.NewStyle()

	var (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")

		headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		cellStyle    = lipgloss.NewStyle().Padding(0, 1)
		oddRowStyle  = cellStyle.Foreground(gray)
		evenRowStyle = cellStyle.Foreground(lightGray)
	)

	var s strings.Builder
	s.WriteString(headerStyle.Render(fmt.Sprintf("Profile: %s", p.Name)))
	s.WriteString("\n")
	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Type", "Source", "Destination")

	for i := range p.Dots {
		dotType := "File"
		if p.Dots[i].IsDir {
			dotType = "Dir"
		}
		t.Row(dotType, p.Dots[i].SrcPath, p.Dots[i].DestPath)
		// fmt.Printf("%s: %s, --> %s\n", dotType, p.Dots[i].SrcPath, p.Dots[i].DestPath)
	}

	s.WriteString(t.Render())
	fmt.Print(s.String())
}

func ReadProfiles() ([]string, error) {
	// TODO: Fix hard coded home path
	file, err := os.Open("/home/daniel/.config/.dottyprofiles")
	if err != nil {
		return []string{}, err
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

		paths = append(paths, strings.Trim(line, " "))
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return paths, nil
}
