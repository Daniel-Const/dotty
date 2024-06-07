package tui

import "github.com/charmbracelet/lipgloss"

var (
    
    title           = lipgloss.NewStyle().
                        Foreground(lipgloss.Color("#5235f2")).
                        PaddingBottom(1).
                        PaddingTop(1).
                        PaddingLeft(1)

    selectHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#e83193"))
    selectDefault   = lipgloss.NewStyle()
)
