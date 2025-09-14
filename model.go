package main

import (
	// "archetype/user"
	"net/http"
	"time"

	m "github.com/SteveMCWin/archetype-common/models"
	tea "github.com/charmbracelet/bubbletea"
)

type AppState uint

const (
	Typing AppState = iota
	Home
	ProfileView
	Settings
	Leaderboard
	About
)

type Tab struct {
	TabName  string
	Contents string
}

type HttpStatus int

type HttpError struct{ error }

func (e HttpError) Error() string { return e.error.Error() }

type Model struct {
	windowWidth  int
	windowHeight int

	online     bool
	state      AppState
	currentTab int
	tabs       []Tab
	theme      ColorTheme
	text       string
	user       m.User
	err        error
}

func NewModel() Model {
	return Model{
		state: Home,
		tabs: []Tab{
			{TabName: "About", Contents: "About page"},
			{TabName: "Settings", Contents: "Settings page"},
			{TabName: "Archetype", Contents: "Here should be some random quote and other stuff"},
			{TabName: "Leaderboard", Contents: "Leaderboard displayed here"},
			{TabName: "Profile", Contents: "Your profile here hehe"},
		},
		currentTab: 2,
		// theme: DefaultTheme,
	}
}

func CheckOnline() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get("https://skilled-gazelle-wildly.ngrok-free.app/api/ping")
	if err != nil {
		return HttpError{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return HttpStatus(res.StatusCode)
}

func LoadLocalData() tea.Msg {
	return "todo"
}

func (m Model) Init() tea.Cmd {
	// perhaps this is where the client should try to get in touch with the server
	return tea.Batch(CheckOnline, tea.SetWindowTitle("Archetype"))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			cmds = append(cmds, CheckOnline)
		case "right", "tab":
			m.currentTab = (m.currentTab+1)%len(m.tabs)
		case "left", "shift+tab":
			m.currentTab = (len(m.tabs)+m.currentTab-1)%len(m.tabs)
		default:
			return m, nil
		}
	case HttpStatus:
		if int(msg) == http.StatusOK {
			m.online = true
		}
	case HttpError:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	return m, tea.Batch(cmds...)
}
