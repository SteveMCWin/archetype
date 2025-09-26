package main

import (
	// "archetype/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mod "github.com/SteveMCWin/archetype-common/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabIndex uint

const (
	About TabIndex = iota
	Settings
	Home
	Leaderboard
	ProfileView
)

type Tab struct {
	TabSymbol     string
	TabName       string
	Contents      string
	StyleContents lipgloss.Style
}

type SupportedTerminals int

const (
	Kitty SupportedTerminals = iota
	Gnome
	Konsole
	Other
)

type HttpStatus int

type HttpError struct{ error }

func (e HttpError) Error() string { return e.error.Error() }

type TestStats struct {
	TypingDuration time.Duration

	Cpm float64
	Wpm float64
	Acc float64
}

type Model struct {
	windowWidth  int
	windowHeight int

	terminal SupportedTerminals
	isOnline bool
	currTab  TabIndex
	tabs     []Tab
	theme    ColorTheme
	user     mod.User
	err      error

	isTyping       bool
	quote          mod.Quote
	typedWord      string
	typedErr       string
	numErr         int
	typedQuote     string
	splitQuote     []string
	typedLen       int
	wordsTyped     int
	quoteCompleted bool
	quoteLoaded    bool
	testBegan      bool
	startTime      time.Time
	stats          TestStats
}

func NewModel() Model {
	return Model{
		tabs: []Tab{ // WARNING: the tabs must be made in the same order as TabIndex definitions. A fix for this would be to make the tabs field a map
			{TabSymbol: "  𝕚  ", TabName: "𝕚 About", Contents: "About page"},
			{TabSymbol: "     ⚙    ", TabName: "⚙ Settings", Contents: "Settings page"},
			{TabSymbol: "    Λ    ", TabName: "Λrchetype", Contents: "", StyleContents: quoteStyle},
			{TabSymbol: "       ♔     ", TabName: "♔ Leaderboard", Contents: "Leaderboard displayed here"},
			{TabSymbol: "    ⚇    ", TabName: "⚇ Profile", Contents: "Your profile here hehe"},
		},
		currTab:  Home,
		theme:    DefaultTheme,
		isTyping: true,
	}
}

func GetQuoteFromServer(length mod.QuoteLen) func() tea.Msg {
	return func() tea.Msg {
		c := &http.Client{
			Timeout: 10 * time.Second,
		}

		res, err := c.Get("https://skilled-gazelle-wildly.ngrok-free.app/api/quote?length=" + strconv.Itoa(int(length)))
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

func SetTerminal() tea.Msg {
	term := Other
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		term = Kitty
	}
	if os.Getenv("KONSOLE_VERSION") != "" {
		term = Konsole
	}
	if os.Getenv("VTE_VERSION") != "" {
		term = Gnome
	}

	return term
}

func ChangeFontSize(term *SupportedTerminals, amount int, pos bool) tea.Cmd {
	var term_name string
	var term_cmd []string
	amount_str := strconv.Itoa(amount)

	switch *term {
	case Kitty:
		term_name = "kitty"
		sign := "+"
		if !pos {
			sign = "-"
		}
		term_cmd = []string{"@", "set-font-size", "--", sign + amount_str}
	}

	c := exec.Command(term_name, term_cmd...) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return nil
	})
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.Sequence(SetTerminal, ChangeFontSize(&m.terminal, 8, true)), // NOTE: Hard coded for testing, the amount should be read from saved user settings
		SetCurrentTheme(m.theme),
		tea.SetWindowTitle("Archetype"),
		GetQuoteFromServer(mod.QUOTE_SHORT),
		tea.ShowCursor,
	}

	return tea.Batch(cmds...) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
}

func ResetTyingData(m *Model) {
	m.numErr = 0
	m.testBegan = false
	m.typedErr = ""
	m.isTyping = true
	m.quoteLoaded = false
	m.wordsTyped = 0
	m.typedWord = ""
	m.typedQuote = ""
	m.quoteCompleted = false
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// log.Println("tea.KeyMsg:", msg.String())
		switch msg.String() {
		case "ctrl+c":
			seq := tea.Sequence(ChangeFontSize(&m.terminal, 0, true), tea.Quit)
			cmds = append(cmds, seq)
		case "right", "tab":
			if !m.isTyping {
				m.currTab = TabIndex((int(m.currTab) + 1) % len(m.tabs))
			}
		case "left", "shift+tab":
			if !m.isTyping {
				m.currTab = TabIndex((len(m.tabs) + int(m.currTab) - 1) % len(m.tabs))
			}
		case "ctrl+r":
			ResetTyingData(&m)
			cmds = append(cmds, GetQuoteFromServer(mod.QUOTE_MEDIUM))
		case "ctrl+up":
			cmds = append(cmds, ChangeFontSize(&m.terminal, 1, true))
		case "ctrl+down":
			cmds = append(cmds, ChangeFontSize(&m.terminal, 1, false))
		case "enter":
			if m.currTab == Home {
				m.isTyping = true
				cmds = append(cmds, tea.ShowCursor)
			}
		case "esc":
			if m.isTyping {
				cmds = append(cmds, tea.HideCursor)
				m.isTyping = false
				// stop the test or something
			}
		default:
			if m.isTyping {
				HandleTyping(&m, msg.String())
			}
		}
	case HttpStatus:
		if int(msg) == http.StatusOK {
			m.isOnline = true
		}
	case HttpError:
		log.Println("ERROR:", msg)
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case mod.Quote:
		m.quoteLoaded = true
		m.quote = msg
		m.splitQuote = strings.Split(m.quote.Quote, " ")
		m.typedLen = len(m.splitQuote[m.wordsTyped])
		m.tabs[Home].Contents = m.quote.Quote
	case SupportedTerminals:
		m.terminal = msg
	}

	return m, tea.Batch(cmds...)
}

func HandleTyping(m *Model, key string) {
	log.Println("key:", key)

	if m.quoteCompleted {
		return
	}

	if !m.testBegan {
		m.testBegan = true
		m.startTime = time.Now()
	}

	switch key {
	case "backspace":
		if m.typedErr != "" {
			m.typedErr = m.typedErr[:max(len(m.typedErr)-1, 0)]
		} else {
			m.typedWord = m.typedWord[:max(len(m.typedWord)-1, 0)]
		}
	case "ctrl+h":
		m.typedErr = ""
		m.typedWord = ""
	case " ":
		if m.typedWord != m.splitQuote[m.wordsTyped] || m.typedErr != "" {
			m.typedErr += " "
			return
		}
		m.wordsTyped++
		m.typedLen += len(m.splitQuote[m.wordsTyped]) + 1
		m.typedQuote += m.typedWord + " "
		m.typedWord = ""
	default:
		if m.typedErr != "" || key != m.splitQuote[m.wordsTyped][len(m.typedWord):min(len(m.typedWord)+1, len(m.splitQuote[m.wordsTyped]))] {
			m.numErr++
			m.typedErr += key
			return
		}
		m.typedWord += key

		if m.wordsTyped+1 == len(m.splitQuote) && m.typedWord == m.splitQuote[m.wordsTyped] {
			m.quoteCompleted = true
			m.isTyping = false
			SetStats(m)
		}
	}

}

func SetStats(m *Model) {
	dur := time.Since(m.startTime)

	cpm := float64(time.Minute.Milliseconds())/float64(dur.Milliseconds()) * float64(len(m.quote.Quote))
	wpm := cpm/4.7
	acc :=  100.0 * (float64(len(m.quote.Quote)-m.numErr) / float64(len(m.quote.Quote)))

	log.Println("Duration:", dur.String())
	log.Println("String len:", len(m.quote.Quote))
	log.Printf("Cpm: %f", cpm)

	m.stats = TestStats{
		TypingDuration: dur,

		Cpm: cpm,
		Wpm: wpm,
		Acc: acc,
	}
}
