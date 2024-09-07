package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
        commands: NewCommandsModel(commands),
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

    switch m.state {
    // Select profile view
    case selectProfile:
        title.WriteString(titleStyle.Render("Dotty · Select a profile"))
        s.WriteString(m.profile.SelectView())

    // Main view
    default:
        title.WriteString(titleStyle.Render("Dotty · Profile: " + m.profile.Profile.Name))
        s.WriteString(
            lipgloss.JoinVertical(
                lipgloss.Top,
                cmdColContainer.Render(m.commands.View()),
                m.profile.ShowView(m.commands.cursor),
            ),
        )
    }
    
    return rootContainer.Render(title.String() + "\n" + s.String())
}
