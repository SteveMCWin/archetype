package main

import (
	"log"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// should check for commandline args
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Unable to run tui:", err)
	}
}

