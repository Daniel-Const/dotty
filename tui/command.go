package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Daniel-Const/dotty/core"
)

const (
    deployCmd int = iota
    loadCmd
)

type runningMsg string 
func runCmd(p *core.Profile, cmd int) tea.Cmd {
    return func() tea.Msg {
        switch cmd {
        case deployCmd:
            p.Deploy()
        case loadCmd:
            p.Load()
        }
        return runningMsg(p.Name)
    } 
}

type Command struct {
    Name    string
    Desc    string
}

// Command UI Model
type CommandsModel struct {
    cursor   int
    running  int
    cmds     []Command
    Profile *core.Profile 
}

func NewCommandsModel(cmds []Command) CommandsModel {
    return CommandsModel{
        cursor:  0,
        running: -1,
        cmds:    cmds,
        Profile: nil,
    }
}

func (m CommandsModel) Init() tea.Cmd {
    return nil
}

func (m CommandsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.cmds)-1 {
                m.cursor++
            }
        case "enter", "":
            m.running = m.cursor
            return m, runCmd(m.Profile, m.cursor)
        }
    }

    return m, nil    
}

func (m CommandsModel) View() string {
    var s strings.Builder
    s.WriteString(fmt.Sprintf("Profile: %s", m.Profile.Name))
    s.WriteString("\n\n")
    for i := range m.cmds {
        if i == m.cursor {
            s.WriteString(selectHighlight.Render("> "+m.cmds[i].Name))
        } else {
            s.WriteString(selectDefault.Render("> "+m.cmds[i].Name))
        }
        s.WriteString("\n")
    }

    if m.running >= 0 {
        s.WriteString(fmt.Sprintf("Running: %s...", m.cmds[m.running].Name))
    }

    return s.String()
}

