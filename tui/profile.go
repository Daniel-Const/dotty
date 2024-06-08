package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Daniel-Const/dotty/core"
)

type submitProfileMsg *core.Profile
type pathErrMsg error
func submitCmd(path string) tea.Cmd {
    return func() tea.Msg {
        // Load profile
        p := core.NewProfile(path)
        if _, err := p.LoadMap(); err != nil {
            return pathErrMsg(err)
        }
        return submitProfileMsg(p)
    }
}

type ProfileModel struct {
    path   textinput.Model
    errMsg string
}

func NewProfileModel() ProfileModel {
    // Initialise path input model
    ti := textinput.New()
    ti.Placeholder = "Profile path"
    ti.Focus()
    ti.CharLimit = 200
    ti.Width = 20
    // TODO: Default profile dir from config?
    ti.SetValue("profiles/daniel-pc")
    return ProfileModel{ti, ""}
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
            return m, submitCmd(m.path.Value())
        }
    case pathErrMsg:
        m.errMsg = msg.Error()
    }

    m.path, cmd = m.path.Update(msg)

    return m, cmd
}

func (m ProfileModel) View() string {
    var s strings.Builder
    s.WriteString("Enter a profile")
    s.WriteString("\n\n")
    s.WriteString(m.path.View())
    s.WriteString("\n")
    s.WriteString(errStyle.Render(m.errMsg))
    return s.String()
}
