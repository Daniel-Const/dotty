package tui

import "github.com/charmbracelet/lipgloss"

// Text
var (
    
    titleStyle      = lipgloss.NewStyle().
                        Foreground(lipgloss.Color("#5235f2")).
                        PaddingTop(1).
                        PaddingLeft(1)

    errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e82d2a"))

    optionHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e83193")).Background(lipgloss.Color("#ffffff"))
    optionDefaultStyle   = lipgloss.NewStyle()
)

// Containers
var (
    rootContainer = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
    cmdColContainer = lipgloss.NewStyle().PaddingRight(2)
)
