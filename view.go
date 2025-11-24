package main

import (
	"fmt"
	"log"

	// "os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

func (m Model) View() tea.View {

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

	log.Println("windowStyle width: ", windowStyle.GetHorizontalFrameSize())
	log.Println("docStyle width: ", docStyle.GetHorizontalFrameSize())
	log.Println("contentTyle width: ", contentStyle.GetHorizontalFrameSize())

	_, err := doc.WriteString(windowStyle.Width(m.windowWidth - windowStyle.GetHorizontalFrameSize()+2).Height(m.windowHeight-windowStyle.GetVerticalFrameSize()).Render(contents))
	if err != nil {
		log.Println("Error displaying window and contents:", err)
	}

	whole_string := docStyle.Render(doc.String())
	if len(m.splitQuote) > 0 {
		position := strings.Index(whole_string, m.splitQuote[0])
		log.Println("Position 1D: ", position)
		log.Println("Position 2D: ", position%m.windowWidth, " ", position/m.windowWidth)
	}

	view := tea.NewView(whole_string)
	view.Cursor = m.cursor
	view.AltScreen = true
	return view
}

func GetHomeContents(m *Model) string {

	// There are two problems: 
	//  - We have to determine the starting position of the text
	//  - We have to determine when the cursor is supposed to go to the next row of text
	// 
	// Ok so the idea is to split the quote into rows, where each row ends with a newline character and then display it as that.
	// Since the text is aligned to the left, the X coordinate of the start of each row is the same
	// So once the startX + currX of the cursor are greater than the length of the text in currRow, set the cursors currX to startX and increment currY
	// I've tried searching the result string that gets rendered to the terminal to get the index of the first word and position the cursor there,
	// but the result string has a bunch of characters that aren't visible and are for style, so it's way off. The only other solution I could
	// think of is to just calculate the position based on the style margins and paddings, which does work for the startX, but the startY is currenlty
	// variable because of the alignment of the content style. This means I will have to handle the alignment myself by calculating the vertical padding
	// for the content based on the window height.

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
		contents = typedStyle.Render("Completed test!!! :D")

		stats_str := fmt.Sprintf("\nWPM: %f\nCPM: %f\nACC: %f\n", m.stats.Wpm, m.stats.Cpm, m.stats.Acc)
		contents += typedStyle.Render(stats_str)
	}

	contents = contentStyle.Width(m.windowWidth-windowStyle.GetHorizontalFrameSize()).Render(contents)

	return contents
}
