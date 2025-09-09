package main

import (
	"archetype/user"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type AppState uint

const (
	Typing AppState = iota
	ProfileView
	Settings
	Leaderboard
	About
)

type Model struct {
	online bool
	state AppState
	theme ColorTheme
	user  user.User
}

func NewModel() Model {
	user, err := user.Load()
	if err != nil {
		// output into file
	}

	return Model{
		state: Typing,
		theme: DefaultTheme,
		user:  user,
	}
}

func (m Model) Init() tea.Cmd {
	// perhaps this is where the client should try to get in touch with the server
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}
