package main

import (
	// "archetype/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	mod "github.com/SteveMCWin/archetype-common/models"
	tea "github.com/charmbracelet/bubbletea"
)

type CurrentTab uint

const (
	About CurrentTab = iota
	Settings
	Home
	Leaderboard
	ProfileView
)

type Tab struct {
	TabSymbol string
	TabName   string
	Contents  string
}

type HttpStatus int

type HttpError struct{ error }

func (e HttpError) Error() string { return e.error.Error() }

type Model struct {
	windowWidth  int
	windowHeight int

	online     bool
	currentTab CurrentTab
	tabs       []Tab
	theme      ColorTheme
	quote      mod.Quote
	user       mod.User
	err        error
}

func NewModel() Model {
	return Model{
		tabs: []Tab{
			{TabSymbol: "  i  ", TabName: "About", Contents: "About page"},
			{TabSymbol: "     ⚙    ", TabName: "⚙ Settings", Contents: "Settings page"},
			{TabSymbol: "    Λ    ", TabName: "Λrchetype", Contents: ""},
			{TabSymbol: "       ♔     ", TabName: "♔ Leaderboard", Contents: "Leaderboard displayed here"},
			{TabSymbol: "    ⚇    ", TabName: "⚇ Profile", Contents: "Your profile here hehe"},
		},
		currentTab: Home,
		theme:      DefaultTheme,
	}
}

func GetQuoteFromServer(length mod.QuoteLen) func() tea.Msg {
	return func() tea.Msg {
		log.Println("Called get quote")
		c := &http.Client{
			Timeout: 10 * time.Second,
		}

		res, err := c.Get("https://skilled-gazelle-wildly.ngrok-free.app/api/quote?length="+strconv.Itoa(int(length)))
		if err != nil {
			return HttpError{err}
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return HttpError{fmt.Errorf("bad status: %s", res.Status)}
		}

		var quote mod.Quote

		if err = json.NewDecoder(res.Body).Decode(&quote); err != nil {
			return HttpError{err}
		}

		log.Println("Quote recieved:", quote)

		return quote
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
	defer res.Body.Close()

	return HttpStatus(res.StatusCode)
}

func LoadLocalData() tea.Msg {
	return "todo"
}

func (m Model) Init() tea.Cmd {
	log.Println()
	log.Println("~~~~~~~~~PROGRAM START~~~~~~~~~")
	log.Println()
	return tea.Batch(GetQuoteFromServer(mod.QUOTE_SHORT), tea.SetWindowTitle("Archetype"), SetCurrentTheme(m.theme)) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
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
			m.currentTab = CurrentTab((int(m.currentTab) + 1) % len(m.tabs))
		case "left", "shift+tab":
			m.currentTab = CurrentTab((len(m.tabs) + int(m.currentTab) - 1) % len(m.tabs))
		case "ctrl+r":
			cmds = append(cmds, GetQuoteFromServer(mod.QUOTE_SHORT))
		default:
			return m, nil
		}
	case HttpStatus:
		if int(msg) == http.StatusOK {
			m.online = true
		}
	case HttpError:
		m.err = msg
		log.Println("ERROR:", m.err)
		return m, nil
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case mod.Quote:
		m.quote = msg
		log.Println("Got the quote")
		m.tabs[Home].Contents = m.quote.Quote
		log.Println("Chagned the contents to:", m.tabs[Home].Contents)
		
	}

	return m, tea.Batch(cmds...)
}
