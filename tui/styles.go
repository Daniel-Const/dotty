package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 44

var (
	green     = lipgloss.Color("#00ff41")
	darkGreen = lipgloss.Color("#007a1f")
	red       = lipgloss.Color("#ff4141")
	grey      = lipgloss.Color("#555555")
	black     = lipgloss.Color("#000000")

	titleStyle    = lipgloss.NewStyle().Foreground(green)
	dimStyle      = lipgloss.NewStyle().Foreground(darkGreen)
	errStyle      = lipgloss.NewStyle().Foreground(red)
	inactiveStyle = lipgloss.NewStyle().Foreground(grey)
	helpStyle     = lipgloss.NewStyle().Foreground(grey)

	cmdActiveStyle   = lipgloss.NewStyle().Foreground(black).Background(green).Bold(true)
	cmdInactiveStyle = lipgloss.NewStyle().Foreground(grey)
)

// retroHeader renders a two-line ASCII box header:
//
//	┌─[ title ]──────────────────────────────┐
//	│  subtitle                              │
//	└────────────────────────────────────────┘
func retroHeader(title, subtitle string) string {
	inner := boxWidth - 2 // chars between │ and │

	titlePart := "─[ " + title + " ]"
	topDashes := strings.Repeat("─", max(0, inner-len([]rune(titlePart))))
	top := "┌" + titlePart + topDashes + "┐"

	midContent := "  " + subtitle
	midPad := strings.Repeat(" ", max(0, inner-len([]rune(midContent))))
	mid := "│" + midContent + midPad + "│"

	bottom := "└" + strings.Repeat("─", inner) + "┘"

	return titleStyle.Render(top) + "\n" +
		dimStyle.Render(mid) + "\n" +
		titleStyle.Render(bottom) + "\n\n"
}
