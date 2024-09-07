package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/*
Load <--
Deploy -->
Edit Map
*/

const (
    deployCmd int = iota
    loadCmd
)

type finishedCmd struct {
    msg string
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
}

func NewCommandsModel(cmds []Command) CommandsModel {
    if len(cmds) == 0 {
        log.Printf("Error creating CommandsModel: No commands given")
        // TODO: Handling errors in constructor?
    }

    return CommandsModel{
        cursor:  0,
        running: -1,
        runMsg: "",
        cmds:    cmds,
    }
}

func (m CommandsModel) Init() tea.Cmd {
    return nil
}

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
        case "enter", "":
            m.running = m.cursor
            return m, nil
        }

    case finishedCmd:
        m.running = -1
        m.runMsg = "Finished running " + msg.msg
    }


    return m, nil    
}

func (m CommandsModel) View() string {
    var s strings.Builder
    s.WriteString("\n\n")
    for i := range m.cmds {
        var optionStyle = optionDefaultStyle
        if i == m.cursor {
            optionStyle = optionHighlightStyle
        }
        s.WriteString(optionStyle.Render(m.cmds[i].Name) + " ")
        // s.WriteString("\n")
    }

    s.WriteString("\n")

    if m.running >= 0 {
        log.Printf("m.running: %d", m.running)
        s.WriteString(fmt.Sprintf("Running: %s...", m.cmds[m.running].Name))
    }

    if m.runMsg != "" {
        s.WriteString(m.runMsg)
    }

    return s.String()
}

