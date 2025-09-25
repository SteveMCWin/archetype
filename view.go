package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isActive := i == int(m.currTab)
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()
		style = style.Border(border)

		if m.isTyping {
			style = style.Faint(true)
		}

		renderedTabs = append(renderedTabs, style.Render(t.TabName))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	gap_size := max(0, m.windowWidth-lipgloss.Width(row)-12)
	gap_l := tabGapLeft.Render(strings.Repeat(" ", gap_size/2))
	gap_r := tabGapRight.Render(strings.Repeat(" ", gap_size/2+gap_size%2))
	row = lipgloss.JoinHorizontal(lipgloss.Top, gap_l, row, gap_r)
	doc.WriteString(row)
	doc.WriteString("\n")

	var contents string

	switch m.currTab {
	case About:
		contents = m.tabs[m.currTab].Contents
	case Settings:
		contents = m.tabs[m.currTab].Contents
	case Home:
		contents = GetHomeContents(&m)
	case Leaderboard:
		contents = m.tabs[m.currTab].Contents
	case ProfileView:
		contents = m.tabs[m.currTab].Contents
	default:
	}

	_, err := doc.WriteString(windowStyle.Width(m.windowWidth - windowStyle.GetHorizontalFrameSize()*3).Height(m.windowHeight-windowStyle.GetVerticalFrameSize()).Render(contents))
	if err != nil {
		log.Println("Error displaying window and contents:", err)
	}

	return docStyle.Render(doc.String())
}

func GetHomeContents(m *Model) string {

	contents := ""

	if m.quoteLoaded && !m.quoteCompleted {
		curr_word := m.splitQuote[m.wordsTyped]

		correct_chars := len(m.typedWord)
		incorrect_chars := len(m.typedErr)
		typed_chars := correct_chars + incorrect_chars

		contents = typedStyle.Render(m.typedQuote) // Already typed words

		contents += typedStyle.Render(curr_word[:min(correct_chars, len(curr_word))]) // Current word - typed correctly
		contents += errorStyle.Render(curr_word[min(correct_chars, len(curr_word)):min(typed_chars, len(curr_word))]) // Current word - typed incorrectly
		contents += quoteStyle.Render(curr_word[min(typed_chars, len(curr_word)):]) // Current word - untyped
		contents += errorStyle.Render(m.typedErr[min(len(curr_word)-correct_chars, incorrect_chars):]) // Current word - overtyped
		
		contents += quoteStyle.Render(m.quote.Quote[m.typedLen:]) // Rest of the quote
	} else if m.quoteCompleted {
		contents = typedStyle.Render("Copleted test!!! :D")

	// TypingDuration time.Duration
	//
	// Cpm float64
	// Wpm float64
	// Acc float64
		stats_str := fmt.Sprintf("\nWPM: %f\nCPM: %f\nACC: %f\n", m.stats.Wpm, m.stats.Cpm, m.stats.Acc)
		contents += typedStyle.Render(stats_str)
	}

	// centerStyle := lipgloss.NewStyle().Margin(0, 8)
	centerStyle := lipgloss.NewStyle().Padding(0, 8)
	contents = centerStyle.Render(contents)

	return contents

}
