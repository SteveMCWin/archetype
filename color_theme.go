package main

import "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
}

var (
	DefaultTheme = ColorTheme {
		Primary: lipgloss.Color("#1e1e2e"),
		Secondary: lipgloss.Color("#6c7086"),
		Accent: lipgloss.Color("#89b4fa"),
	}
)
