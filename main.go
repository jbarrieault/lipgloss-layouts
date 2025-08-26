package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	w, h       int
	headerText string
	bodyText   string
	footerText string
}

func initialModel() model {
	return model{
		headerText: "Header Title",
		bodyText: `This is the body content that can either grow or shink to fill available space.

It can contain multiple lines of text, and may be constrained by the terminal height to prevent header and footer from overflowing.

Or, overflowing may be desired.

This various approaches to layout management in for TUI applications using lipgloss.

Try lipgloss, you'll definitely love it!

I promise.`,
		footerText: "Footer — press q to quit",
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

var (
	pink   = lipgloss.Color("#8b5a8b")
	green  = lipgloss.Color("#4a7c4a")
	blue   = lipgloss.Color("#4a5f8b")
	orange = lipgloss.Color("#b8860b")

	container = lipgloss.NewStyle().Padding(1).Background(pink)

	header = lipgloss.NewStyle().
		Bold(true).
		Background(green).
		Foreground(lipgloss.Color("#f0f0f0")).
		Align(lipgloss.Center)

	body = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("#f0f0f0")).
		Background(blue)

	footer = lipgloss.NewStyle().
		Background(orange).
		Foreground(lipgloss.Color("#f0f0f0")).
		Align(lipgloss.Center)
)

func (m model) View() string {
	if m.w == 0 || m.h == 0 {
		return "loading…"
	}

	screenWidth := m.w
	screenHeight := m.h

	containerWidth, containerHeight := container.GetFrameSize()
	innerW := screenWidth - containerWidth
	innerH := screenHeight - containerHeight

	hdr := header.Width(innerW).Render(m.headerText)
	ftr := footer.Width(innerW).Render(m.footerText)

	// Get body style's frame size to determine how much margin, border, and padding it takes up
	bodyContainerWidth, bodyContainerHeight := body.GetFrameSize()
	innerBodyWidth := innerW - bodyContainerWidth
	innerBodyHeight := innerH - bodyContainerHeight

	// Account for header, footer
	used := lipgloss.Height(hdr) + lipgloss.Height(ftr)
	availableBodyH := innerBodyHeight - used

	// A few different layout styling options below...

	// 1. Shrink Body
	// - body height: only use as much as needed for its content
	// - footer: sticks to bottom of body (not the terminal screen)
	// - content overflow: body & footer overflow off bottom of terminal screen when content doesn't fit
	//
	// Solid option, so long as the application implements a
	// "Your terminal is too small" screen when content height > terminal height,
	// if your footer contains useful info such as keyboard shortcuts.
	// bodyStyle := body.Width(innerBodyWidth)

	// 2. Grow Body
	// - body height: grows (inside border) to to fill available screen space
	// - footer: appears at bottom of screen (because body takes up available space)
	// - content overflow: body & footer overflow off bottom of terminal screen when content doesn't fit
	//
	// This option is very similar to 1.
	// bodyStyle := body.Width(innerBodyWidth).Height(availableBodyH)

	// 3. Shrink Body + no footer overflow
	// - body height: only use as much as needed for its content
	// - footer: appears at bottom of screen (because body takes up available space)
	// - content overflow: body is clipped at its bottom if it's content doesn't fit,
	//     but the footer is not clipped and remained visible at the bottom of the screen.
	//
	// The only weird thing about this one is that when the body content is clipped,
	// the sticky footer covers the body's bottom border.
	// In that case, the user may not realize the body content is being truncated.
	// But, if body has a border, that being covered by the footer may give enough
	// visual indication that the content is truncated.
	// bodyStyle := body.Width(innerBodyWidth).MaxHeight(availableBodyH)

	// 4. Grow Body + no footer overflow
	// I don't fully understand why this works, but it does.
	bodyStyle := body.Width(innerBodyWidth).Height(availableBodyH).MaxHeight(innerBodyHeight)

	// log.Printf("height: %d, width: %d, innerW: %d, innerH: %d, containerWidth: %d, containerHeight: %d", screenHeight, screenWidth, innerW, innerH, containerWidth, containerHeight)

	bdy := bodyStyle.Render(m.bodyText)

	content := lipgloss.JoinVertical(lipgloss.Left, hdr, bdy, ftr)

	// setting MaxHeight on the main container is key to making
	// the top of the container's content fixed to the top of the terminal screen,
	// even when the terminal size is too small to fit all content.
	return container.Width(screenWidth).MaxHeight(screenHeight).Render(content)
}

func main() {
	f, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
