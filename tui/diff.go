package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Daniel-Const/dotty/core"
)

type exitDiffMsg struct{}

type FileDiff struct {
	Name    string
	Changed bool
	IsDir   bool
	Lines   []core.DiffLine
}

type DiffModel struct {
	profile *core.Profile
	diffs   []FileDiff
	cursor  int
}

func NewDiffModel(profile *core.Profile) DiffModel {
	return DiffModel{
		profile: profile,
		diffs:   computeDiffs(profile),
		cursor:  0,
	}
}

func computeDiffs(profile *core.Profile) []FileDiff {
	var result []FileDiff
	for _, dot := range profile.Dots {
		name := strings.TrimPrefix(strings.TrimPrefix(dot.SrcPath, profile.Location), "/")
		if dot.IsDir {
			result = append(result, FileDiff{Name: name, IsDir: true})
			continue
		}
		lines, err := core.DiffFile(dot.SrcPath, dot.DestPath)
		if err != nil {
			result = append(result, FileDiff{Name: name})
			continue
		}
		changed := false
		for _, l := range lines {
			if l.Kind != "same" {
				changed = true
				break
			}
		}
		result = append(result, FileDiff{Name: name, Changed: changed, Lines: lines})
	}
	return result
}

func (m DiffModel) Init() tea.Cmd { return nil }

func (m DiffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.diffs)-1 {
				m.cursor++
			}
		case "q", "esc":
			return m, func() tea.Msg { return exitDiffMsg{} }
		}
	}
	return m, nil
}

func (m DiffModel) View() string {
	var s strings.Builder

	for i, fd := range m.diffs {
		name := fd.Name
		pad := strings.Repeat(" ", max(0, 24-len([]rune(name))))
		cursor := "    "
		if i == m.cursor {
			cursor = "  ▶ "
		}

		if i == m.cursor {
			if fd.Changed {
				s.WriteString(errStyle.Render(cursor+name+pad+"[ CHANGED ]") + "\n")
			} else {
				s.WriteString(titleStyle.Render(cursor+name+pad+"[ IN SYNC ]") + "\n")
			}
		} else {
			if fd.Changed {
				s.WriteString(inactiveStyle.Render(cursor+name+pad) + errStyle.Render("[ CHANGED ]") + "\n")
			} else {
				s.WriteString(dimStyle.Render(cursor+name+pad+"[ IN SYNC ]") + "\n")
			}
		}
	}

	if len(m.diffs) > 0 {
		fd := m.diffs[m.cursor]
		divider := "  ┄┄┄┄┄┄ " + fd.Name + " " + strings.Repeat("┄", max(0, 30-len([]rune(fd.Name))))
		s.WriteString("\n" + dimStyle.Render(divider) + "\n\n")

		switch {
		case fd.IsDir:
			s.WriteString(inactiveStyle.Render("  (directory — not diffable)") + "\n")
		case !fd.Changed:
			s.WriteString(dimStyle.Render("  ✓ in sync") + "\n")
		default:
			for _, line := range fd.Lines {
				switch line.Kind {
				case "add":
					s.WriteString(titleStyle.Render("  + "+line.Text) + "\n")
				case "remove":
					s.WriteString(errStyle.Render("  - "+line.Text) + "\n")
				case "same":
					s.WriteString(inactiveStyle.Render("    "+line.Text) + "\n")
				}
			}
		}
	}

	s.WriteString("\n" + helpStyle.Render("  ↑↓ navigate files   q back to menu"))
	return s.String()
}
