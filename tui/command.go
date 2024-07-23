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

type finishedCmd struct {
    msg string
}

func runCmd(p *core.Profile, cmd int) tea.Cmd {
    return func() tea.Msg {
        cmdStr := ""
        switch cmd {
        case deployCmd:
            cmdStr = "Deploy"
            p.Deploy()
        case loadCmd:
            cmdStr = "Load"
            p.Load()
        }

        return finishedCmd{cmdStr}
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
    runMsg   string 
    cmds     []Command
    Profile *core.Profile 
}

func NewCommandsModel(cmds []Command) CommandsModel {
    return CommandsModel{
        cursor:  0,
        running: -1,
        runMsg: "",
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
    case finishedCmd:
        m.running = -1
        m.runMsg = "Finished running " + msg.msg
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

    s.WriteString("\n")

    if m.running >= 0 {
        s.WriteString(fmt.Sprintf("Running: %s...", m.cmds[m.running].Name))
    }

    if m.runMsg != "" {
        s.WriteString(m.runMsg)
    }

    return s.String()
}

