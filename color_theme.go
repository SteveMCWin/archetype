package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ColorTheme struct {
	Primary    lipgloss.AdaptiveColor
	Secondary  lipgloss.AdaptiveColor
	Accent     lipgloss.AdaptiveColor
	TextError  lipgloss.AdaptiveColor
	TextTyped  lipgloss.AdaptiveColor
	TextUnyped lipgloss.AdaptiveColor
}

var (
	DefaultTheme = ColorTheme{
		Primary:    lipgloss.AdaptiveColor{Dark: "#1e1e2e", Light: "#6c7086"},
		Secondary:  lipgloss.AdaptiveColor{Dark: "#6c7086", Light: "#acb0be"},
		Accent:     lipgloss.AdaptiveColor{Dark: "#89b4fa", Light: "#1e66f5"},
		TextError:  lipgloss.AdaptiveColor{Dark: "#dd8888", Light: "#dd8888"},
		TextTyped:  lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		TextUnyped: lipgloss.AdaptiveColor{Dark: "#aaaaaa", Light: "#444444"},
	}
)

var (
	inactiveTabBorder = lipgloss.Border{Bottom: "─", BottomLeft: "─", BottomRight: "─"}
	activeTabBorder   = lipgloss.Border{Top: "─", Bottom: " ", Left: "│", Right: "│", TopLeft: "╭", TopRight: "╮", BottomLeft: "┘", BottomRight: "└"}
	tabGapBorderLeft  = lipgloss.Border{Bottom: "─", BottomLeft: "╭", BottomRight: "─"}
	tabGapBorderRight = lipgloss.Border{Bottom: "─", BottomLeft: "─", BottomRight: "╮"}
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	tabGapLeft        = inactiveTabStyle.Border(tabGapBorderLeft, true)
	tabGapRight       = inactiveTabStyle.Border(tabGapBorderRight, true)
	windowStyle       = lipgloss.NewStyle().Padding(2, 0).Align(lipgloss.Left, lipgloss.Center).Border(lipgloss.RoundedBorder()).UnsetBorderTop()
	quoteStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextUnyped)
	typedStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextTyped)
	errorStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextError)
)

func SetCurrentTheme(t ColorTheme) func() tea.Msg {
	return func() tea.Msg {
		inactiveTabStyle = inactiveTabStyle.BorderForeground(t.Accent)
		activeTabStyle = activeTabStyle.BorderForeground(t.Accent)
		tabGapLeft = tabGapLeft.BorderForeground(t.Accent)
		tabGapRight = tabGapRight.BorderForeground(t.Accent)
		windowStyle = windowStyle.BorderForeground(t.Accent)
		return nil
	}
}
