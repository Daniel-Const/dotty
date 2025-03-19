package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
}

func NewProfileModel() ProfileModel {
	// Initialise path input model
	ti := textinput.New()
	ti.Placeholder = "Profile path"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 60

	defaultPath := ""
	home, err := os.UserHomeDir()
	if err == nil {
		defaultPath = filepath.Join(home, "/")
	}

	ti.SetValue(defaultPath)
	profiles, err := core.ReadProfileList()
	if err != nil {
		fmt.Println(err)
	}

	return ProfileModel{ti, "", nil, 0, len(profiles), profiles}
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
			// TODO: Save new profile to .config/.dottyprofiles
			// Load new profile from the path
			path := ""
			if m.CursorAtEnd() {
				path = m.path.Value()
			} else {
				path = m.profiles[m.cursor]
			}
			p := core.NewProfile(path)
			if _, err := p.LoadMap(); err != nil {
				return m, submitCmd(err)
			}
			m.Profile = p
			return m, submitCmd(nil)
		case "down":
			if m.cursor < m.maxCursor {
				m.cursor += 1
			}
		case "up":
			if m.cursor > 0 {
				m.cursor -= 1
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

func (m ProfileModel) ShowView(direction int) string {
	srcCol := strings.Builder{}
	dirCol := strings.Builder{}
	destCol := strings.Builder{}
	dirChar := "=>"
	if direction == 1 {
		dirChar = "<="
	}
	for _, dot := range m.GetDots() {
		src := strings.ReplaceAll(dot.SrcPath, m.Profile.Location+"/", "")
		srcCol.WriteString(src + "\n")
		dirCol.WriteString("  " + dirChar + "  \n")
		destCol.WriteString(dot.DestPath + "\n")
		// s.WriteString(src + " => " + dot.DestPath + "\n")
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Right,
		srcCol.String(),
		dirCol.String(),
		destCol.String(),
	)
}

func (m ProfileModel) SelectView() string {
	var s strings.Builder

	s.WriteString("\n")

	for i, p := range m.profiles {
		if i == m.cursor {
			s.WriteString(fmt.Sprintf("> %s\n", p))
		} else {
			s.WriteString(fmt.Sprintf("%s\n", p))
		}
	}

	s.WriteString("\n")

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render("New Profile: "))
	if m.cursor >= m.maxCursor {
		s.WriteString(m.path.View())
	}
	s.WriteString("\n\n")
	s.WriteString(errStyle.Render(m.errMsg))

	return s.String()
}

func (m ProfileModel) View() string {
	return ""
}
