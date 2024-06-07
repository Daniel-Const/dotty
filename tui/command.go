package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Daniel-Const/dotty/core"
)

// TODO: Better way to map cursor to command?
const (
    deployCmd int = iota
    loadCmd
)

type runningMsg string 
func profileCmd(cmd int) tea.Cmd {
    return func() tea.Msg {
        // TODO: Obtain from the TUI
        profilePath := "./profiles/daniel-pc"
        p := core.NewProfile(profilePath).LoadMap()
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
    Profile *core.Profile 
}

func NewCommand(name string, desc string) Command {
    c := Command{name, desc, nil}
    return c
}
