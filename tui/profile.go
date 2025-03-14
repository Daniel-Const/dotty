package tui

import (
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
	path    textinput.Model
	errMsg  string
	Profile *core.Profile
}

func NewProfileModel() ProfileModel {
	// Initialise path input model
	ti := textinput.New()
	ti.Placeholder = "Profile path"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 60

	// TODO: Improve profile loading (search for map files?)
	defaultPath := ""
	home, err := os.UserHomeDir()
	if err == nil {
		defaultPath = filepath.Join(home, "/dotfiles"+"/arch-desktop") // TODO: + arch-desktop for easy testing (remove)
	}

	ti.SetValue(defaultPath)
	return ProfileModel{ti, "", nil}
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
			// Load new profile from the path
			p := core.NewProfile(m.path.Value())
			if _, err := p.LoadMap(); err != nil {
				return m, submitCmd(err)
			}
			m.Profile = p
			return m, submitCmd(nil)
		}
	case pathErrMsg:
		m.errMsg = msg.Error()
	}

	m.path, cmd = m.path.Update(msg)

	return m, cmd
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

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render("Profile path: "))
	s.WriteString(m.path.View())
	s.WriteString("\n\n")
	s.WriteString(errStyle.Render(m.errMsg))

	return s.String()
}

func (m ProfileModel) View() string {
	return ""
}
