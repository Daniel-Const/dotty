package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	deployCmd int = iota
	loadCmd
	diffCmd
)

type triggerCmdMsg struct{ cmd int }

func triggerCmd(cmd int) tea.Cmd {
	return func() tea.Msg { return triggerCmdMsg{cmd} }
}

type finishedCmd struct{ msg string }

type Command struct {
	Name string
	Desc string
}

type CommandsModel struct {
	cursor  int
	running int
	cmds    []Command
}

func NewCommandsModel(cmds []Command) CommandsModel {
	return CommandsModel{cursor: 0, running: -1, cmds: cmds}
}

func (m CommandsModel) Init() tea.Cmd { return nil }

func (m CommandsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l":
			if m.cursor < len(m.cmds)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.running == -1 {
				m.running = m.cursor
				return m, triggerCmd(m.running)
			}
		}
	case finishedCmd:
		m.running = -1
	}
	return m, nil
}

func (m CommandsModel) CommandSelectView() string {
	var s strings.Builder
	s.WriteString("\n")
	for i, cmd := range m.cmds {
		label := "[ " + strings.ToUpper(cmd.Name) + " ]"
		if i == m.cursor {
			s.WriteString("  " + cmdActiveStyle.Render(label))
		} else {
			s.WriteString("  " + cmdInactiveStyle.Render(label))
		}
	}
	s.WriteString("\n")
	return s.String()
}

func (m CommandsModel) RunningName() string {
	if m.running >= 0 && m.running < len(m.cmds) {
		return strings.ToUpper(m.cmds[m.running].Name)
	}
	return "COMMAND"
}

func (m CommandsModel) View() string { return "" }
