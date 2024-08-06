package tui

import "github.com/charmbracelet/lipgloss"

// Text
var (
    
    titleStyle      = lipgloss.NewStyle().
                        Foreground(lipgloss.Color("#5235f2")).
                        PaddingBottom(1).
                        PaddingTop(1).
                        PaddingLeft(1)

    selectHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#e83193"))
    selectDefault   = lipgloss.NewStyle()
    errStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#e82d2a"))
)

// Containers
var (
    rootContainer = lipgloss.NewStyle().
                        BorderStyle(lipgloss.NormalBorder())
)
