package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
    selectProfile int = iota
    selectCommand
)

// Main bubbletea model for the app
type Model struct {
    profile  ProfileModel
    commands CommandsModel
    state    int
}

func NewModel(commands []Command) Model {
    return Model{
        profile:  NewProfileModel(),
        state:    selectProfile,
    } 
}

func (m Model) Init() tea.Cmd {
    var cmds []tea.Cmd
    cmds = append(cmds, m.profile.Init())
    return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        }
    case submitProfileMsg:
        m.state++
    }

    cmd := m.updateBubbles(msg)
    return m, cmd
}

func (m *Model) updateBubbles(msg tea.Msg) tea.Cmd {
    var cmds []tea.Cmd
    switch m.state {
    case selectProfile:
        model, cmd := m.profile.Update(msg)
        if p, ok := model.(ProfileModel); ok {
           m.profile = p
        }
        cmds = append(cmds, cmd)
    case selectCommand:
        model, cmd := m.commands.Update(msg)
        if c, ok := model.(CommandsModel); ok {
            m.commands = c
        } 
        cmds = append(cmds, cmd)
    }

    return tea.Batch(cmds...)
}

func (m Model) View() string {
    var s, title strings.Builder

    title.WriteString(titleStyle.Render("Dotty") + "\n")

    switch m.state {
    case selectProfile:
        s.WriteString(m.profile.View())
        return rootContainer.Render(title.String() + s.String())
    }

    title.WriteString("Profile: " + m.profile.Profile.Name)
    s.WriteString(m.profile.ViewMap())

    // TODO: Render buttons (deploy, load, etc...)
    // Implement command cursor for highlight + reverse arrows in ViewMap

    return rootContainer.Render(title.String() + "\n\n" + s.String())
}
