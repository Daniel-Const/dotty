package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Daniel-Const/dotty/core"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	selectProfile int = iota
	selectCommand
	runningCommand
	viewingResult
	viewingDiff
)

type returnToMenuMsg struct{}

func returnAfterDelay() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return returnToMenuMsg{}
	})
}

func runCmd(cmd int, p *core.Profile) tea.Cmd {
	return func() tea.Msg {
		var count int
		var name string
		switch cmd {
		case deployCmd:
			count, _ = p.Deploy()
			name = "Deployed"
		case loadCmd:
			count, _ = p.Load()
			name = "Loaded"
		}
		return finishedCmd{msg: fmt.Sprintf("✓ %s %d files", name, count)}
	}
}

type Model struct {
	profile   ProfileModel
	commands  CommandsModel
	diff      DiffModel
	state     int
	resultMsg string
}

func NewModel(commands []Command, config *core.DottyConfig) Model {
	return Model{
		profile:  NewProfileModel(config),
		commands: NewCommandsModel(commands),
		state:    selectProfile,
	}
}

func (m Model) Init() tea.Cmd {
	return m.profile.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			switch m.state {
			case viewingResult:
				m.state = selectCommand
				return m, nil
			case viewingDiff:
				// DiffModel handles q and returns exitDiffMsg
			default:
				return m, tea.Quit
			}
		case "enter", "esc":
			if m.state == viewingResult {
				m.state = selectCommand
				return m, nil
			}
		}
	case submitProfileMsg:
		m.state = selectCommand
		return m, nil
	case triggerCmdMsg:
		if msg.cmd == diffCmd {
			m.diff = NewDiffModel(m.profile.Profile)
			m.state = viewingDiff
			return m, nil
		}
		m.state = runningCommand
		return m, runCmd(msg.cmd, m.profile.Profile)
	case finishedCmd:
		m.commands.running = -1
		m.resultMsg = msg.msg
		m.state = viewingResult
		return m, returnAfterDelay()
	case returnToMenuMsg:
		if m.state == viewingResult {
			m.state = selectCommand
		}
		return m, nil
	case exitDiffMsg:
		m.state = selectCommand
		return m, nil
	}

	return m, m.updateBubbles(msg)
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
	case selectCommand:
		model, cmd := m.commands.Update(msg)
		if c, ok := model.(CommandsModel); ok {
			m.commands = c
		}
		cmds = append(cmds, cmd)
	case viewingDiff:
		model, cmd := m.diff.Update(msg)
		if d, ok := model.(DiffModel); ok {
			m.diff = d
		}
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var s strings.Builder

	switch m.state {
	case selectProfile:
		s.WriteString(retroHeader("DOTTY", "dotfiles manager"))
		s.WriteString(m.profile.SelectView())

	case selectCommand:
		s.WriteString(retroHeader("DOTTY", "profile: "+m.profile.Profile.Name))
		s.WriteString(m.commands.CommandSelectView())
		s.WriteString(dimStyle.Render("\n  ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄") + "\n\n")
		s.WriteString(m.profile.ShowView())
		s.WriteString("\n" + helpStyle.Render("  ← → navigate   enter run   q quit"))

	case runningCommand:
		s.WriteString(retroHeader("DOTTY", "profile: "+m.profile.Profile.Name))
		s.WriteString("\n  " + titleStyle.Render("Running: "+m.commands.RunningName()+"...") + "\n")

	case viewingResult:
		s.WriteString(retroHeader("DOTTY", "profile: "+m.profile.Profile.Name))
		s.WriteString(resultView(m.resultMsg))

	case viewingDiff:
		s.WriteString(retroHeader("DIFF", "profile: "+m.profile.Profile.Name))
		s.WriteString(m.diff.View())
	}

	return s.String()
}

const resultInner = 40

func resultView(msg string) string {
	content := "  " + msg
	spaces := strings.Repeat(" ", max(0, resultInner-len([]rune(content))))
	topDashes := strings.Repeat("─", max(0, resultInner-11)) // after "─[ RESULT ]"

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(dimStyle.Render("  ┌─[ RESULT ]"+topDashes+"┐") + "\n")
	s.WriteString(dimStyle.Render("  │"+strings.Repeat(" ", resultInner)+"│") + "\n")
	s.WriteString(dimStyle.Render("  │") + titleStyle.Render(content) + spaces + dimStyle.Render("│") + "\n")
	s.WriteString(dimStyle.Render("  │"+strings.Repeat(" ", resultInner)+"│") + "\n")
	s.WriteString(dimStyle.Render("  └"+strings.Repeat("─", resultInner)+"┘") + "\n")
	s.WriteString("\n" + helpStyle.Render("  returning to menu..."))
	return s.String()
}
