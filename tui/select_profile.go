package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Daniel-Const/dotty/core"
)

type submitProfileMsg struct{}
type pathErrMsg error

func submitCmd(err error) tea.Cmd {
	return func() tea.Msg {
		if err != nil {
			return pathErrMsg(err)
		}
		return submitProfileMsg{}
	}
}

type ProfileModel struct {
	path      textinput.Model
	errMsg    string
	Profile   *core.Profile
	cursor    int
	maxCursor int
	profiles  []string
	config    *core.DottyConfig
}

func NewProfileModel(config *core.DottyConfig) ProfileModel {
	ti := textinput.New()
	ti.Placeholder = "Profile path"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 38

	defaultPath := ""
	home, err := os.UserHomeDir()
	if err == nil {
		defaultPath = filepath.Join(home, "/")
	}
	ti.SetValue(defaultPath)

	profiles := config.Profiles
	return ProfileModel{ti, "", nil, 0, len(profiles), profiles, config}
}

func (m ProfileModel) GetDots() []*core.Dot {
	if m.Profile != nil {
		return m.Profile.Dots
	}
	return nil
}

func (m ProfileModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			path := ""
			isManual := m.CursorAtEnd()
			if isManual {
				path = m.path.Value()
			} else {
				path = m.profiles[m.cursor]
			}
			p := core.NewProfile(path)
			if _, err := p.LoadMap(); err != nil {
				return m, submitCmd(err)
			}
			m.Profile = p
			if isManual {
				m.config.AddProfile(path)
				m.config.Save()
			}
			return m, submitCmd(nil)
		case "down":
			if m.cursor < m.maxCursor {
				m.cursor++
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}

	case pathErrMsg:
		m.errMsg = msg.Error()
	}

	m.path, cmd = m.path.Update(msg)
	return m, cmd
}

func (m ProfileModel) CursorAtEnd() bool {
	return m.cursor >= m.maxCursor
}

func (m ProfileModel) ShowView() string {
	var s strings.Builder
	for _, dot := range m.GetDots() {
		src := strings.TrimPrefix(strings.TrimPrefix(dot.SrcPath, m.Profile.Location), "/")
		dots := strings.Repeat("·", max(0, 24-len([]rune(src))))
		s.WriteString(inactiveStyle.Render("  "+src+" ") + dimStyle.Render(dots+" ") + inactiveStyle.Render(dot.DestPath) + "\n")
	}
	return s.String()
}

func (m ProfileModel) SelectView() string {
	var s strings.Builder
	s.WriteString("\n")

	for i, p := range m.profiles {
		if i == m.cursor {
			s.WriteString(errStyle.Render("  ▶ ") + titleStyle.Render(p) + "\n")
		} else {
			s.WriteString(inactiveStyle.Render("    "+p) + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(dimStyle.Render("  New Profile: "))
	if m.CursorAtEnd() {
		s.WriteString(m.path.View())
	}
	s.WriteString("\n\n")

	if m.errMsg != "" {
		s.WriteString(errStyle.Render("  "+m.errMsg) + "\n")
	}

	s.WriteString("\n" + helpStyle.Render(fmt.Sprintf("  ↑↓ navigate   enter select   q quit")))
	return s.String()
}

func (m ProfileModel) View() string { return "" }
