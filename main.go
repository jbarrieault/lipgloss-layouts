package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	w, h             int
	headerText       string
	bodyText         string
	footerText       string
	leftPaneContent  string
	rightPaneContent string
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
		leftPaneContent:  "[Left Pane]\n\nNot a lot of content here.",
		rightPaneContent: "[Right Pane]\n\nMore content here\nthan in the left pane so it will likely require more vertical space.",
		footerText:       "Footer — press q to quit",
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
	red    = lipgloss.Color("#c91b12")
	purple = lipgloss.Color("#8f12c9")

	containerStyle = lipgloss.NewStyle().Padding(1).Background(pink)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Background(green).
			Foreground(lipgloss.Color("#f0f0f0")).
			Align(lipgloss.Center)

	bodyStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderBackground(blue).
			Foreground(lipgloss.Color("#f0f0f0")).
			Background(blue)

	leftPaneStyle = lipgloss.NewStyle().
			Background(red).
			Foreground(lipgloss.Color("#f0f0f0")).
			Align(lipgloss.Left)

	rightPaneStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("#f0f0f0")).
			Align(lipgloss.Right)

	footerStyle = lipgloss.NewStyle().
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

	containerWidth, containerHeight := containerStyle.GetFrameSize()
	innerW := screenWidth - containerWidth
	innerH := screenHeight - containerHeight

	hdr := headerStyle.Width(innerW).Render(m.headerText)
	ftr := footerStyle.Width(innerW).Render(m.footerText)

	// Get body style's frame size to determine how much margin, border, and padding it takes up
	bodyContainerWidth, bodyContainerHeight := bodyStyle.GetFrameSize()
	innerBodyWidth := innerW - bodyContainerWidth
	innerBodyHeight := innerH - bodyContainerHeight

	// Account for header, footer
	used := lipgloss.Height(hdr) + lipgloss.Height(ftr)
	availableBodyH := innerBodyHeight - used

	// Grow body pane heights + no footer overflow
	// I don't fully understand why this combination of Height and MaxHeight works,
	// but it does correctly grow the body pane heights to fill the available space,
	// while also making the footer "sticky".
	leftPane := leftPaneStyle.Width(innerBodyWidth / 2).Height(availableBodyH).MaxHeight(availableBodyH).Render(m.leftPaneContent)
	rightPane := rightPaneStyle.Width(innerBodyWidth / 2).Height(availableBodyH).MaxHeight(availableBodyH).Render(m.rightPaneContent)
	bodyContent := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	bdy := bodyStyle.Render(bodyContent)

	content := lipgloss.JoinVertical(lipgloss.Left, hdr, bdy, ftr)

	// log.Printf("height: %d, width: %d, innerW: %d, innerH: %d, containerWidth: %d, containerHeight: %d", screenHeight, screenWidth, innerW, innerH, containerWidth, containerHeight)

	// setting MaxHeight on the main container is key to making
	// the top of the container's content fixed to the top of the terminal screen,
	// even when the terminal size is too small to fit all content.
	return containerStyle.Width(screenWidth).MaxHeight(screenHeight).Render(content)
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
