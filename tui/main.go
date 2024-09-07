package tui

import (
	"log"
	"strings"
	"time"

	"github.com/Daniel-Const/dotty/core"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	selectProfile int = iota
	selectCommand     // change to profile
	runningCommand
)

func runCmd(cmd int, p *core.Profile) tea.Cmd {
	return func() tea.Msg {
		switch cmd {
		case deployCmd:
			p.Deploy()
		case loadCmd:
			p.Load()
		}

		// Sleep for user experience
		time.Sleep(2 * time.Second)

		return finishedCmd{}
	}

}

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
		case "enter":
			// Return to profile view
			if m.state == runningCommand {
				m.state = selectCommand
				return m, nil
			}
		}
	case submitProfileMsg:
		m.state = selectCommand
	case triggerCmdMsg:
		log.Println("Run command message")
		m.state = runningCommand
		return m, runCmd(msg.cmd, m.profile.Profile)
	}

	cmd := m.updateBubbles(msg)
	return m, cmd
}

func (m *Model) updateBubbles(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch m.state {
	case selectProfile:
		model, cmd := m.profile.Update(msg)
		if pm, ok := model.(ProfileModel); ok {
			m.profile = pm
		}
		cmds = append(cmds, cmd)
	case selectCommand, runningCommand:
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

	// Main profile view
	default:
		title.WriteString(titleStyle.Render("Dotty · Profile: " + m.profile.Profile.Name))
		s.WriteString(
			lipgloss.JoinVertical(
				lipgloss.Top,
				cmdColContainer.Render(m.commands.CommandSelectView()),
				m.profile.ShowView(m.commands.cursor),
			),
		)
	}

	return rootContainer.Render(title.String(), "\n", s.String(), "\n", m.commands.View())
}
