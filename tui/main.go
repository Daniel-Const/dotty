package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)




// Main bubbletea model for the app
type Model struct {
    commands []Command
    cursor   int
    running  int
}

func NewModel(commands []Command) Model {
    m := Model{commands, 0, -1} 
    return m
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    // Key press
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.commands)-1 {
                m.cursor++
            }
        case "enter", "":
            // TODO: Implement actually run the command
            m.running = m.cursor
            return m, profileCmd(m.cursor)
        }
    }
    return m, nil
}

func (m Model) View() string {
    var s strings.Builder

    s.WriteString(title.Render("Dotty"))
    s.WriteRune('\n')

    // Select command
    for i := range m.commands {
        if i == m.cursor {
            s.WriteString(selectHighlight.Render("> "+m.commands[i].Name))
        } else {
            s.WriteString(selectDefault.Render("> "+m.commands[i].Name))
        }
        s.WriteRune('\n')
    }

    // Command is running...
    if m.running >= 0 {
        s.WriteString(fmt.Sprintf("Running: %s...", m.commands[m.running].Name))
    }

    return s.String() 
}
