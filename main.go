package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tick message for animation
type tickMsg time.Time

// Star represents each element on screen
type Star struct {
	X       int
	Y       int
	Visible bool
	Symbol  string
}

// Number of stars/elements
const StarCount = 40

// Theme struct holds the colors and symbols
type Theme struct {
	Name       string
	Background string
	Foreground string
	Symbols    []string
}

// List of themes
var themes = []Theme{
	{
		Name:       "Starry Sky",
		Background: "#0B1020",
		Foreground: "#FAFAFA",
		Symbols:    []string{"✨"},
	},
	{
		Name:       "Flowers",
		Background: "#1B2F2B",
		Foreground: "#FFD3E0",
		Symbols:    []string{"🌸", "🌼", "🌺"},
	},
	{
		Name:       "Ocean",
		Background: "#0E3B5F",
		Foreground: "#A3DFF7",
		Symbols:    []string{"🌊", "💧"},
	},
}

// Main program model
type model struct {
	stars  []Star
	theme  int
	width  int
	height int
}

// Lipgloss style for a theme
func styleForTheme(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Foreground)).
		Background(lipgloss.Color(t.Background)).
		Padding(1, 2)
}

// Initialize stars randomly for a theme
func initStars(t Theme) []Star {
	stars := make([]Star, StarCount)
	for i := 0; i < StarCount; i++ {
		stars[i] = Star{
			X:       rand.Intn(20), // fixed columns for simplicity
			Y:       rand.Intn(10), // fixed rows for simplicity
			Visible: rand.Intn(2) == 0,
			Symbol:  t.Symbols[rand.Intn(len(t.Symbols))],
		}
	}
	return stars
}

// Initialize program
func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update function handles ticks and keypresses
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Toggle some stars on/off
		for i := range m.stars {
			if rand.Intn(100) < 30 {
				m.stars[i].Visible = !m.stars[i].Visible
			}
		}
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.theme = 0
			m.stars = initStars(themes[m.theme])
		case "2":
			m.theme = 1
			m.stars = initStars(themes[m.theme])
		case "3":
			m.theme = 2
			m.stars = initStars(themes[m.theme])
		}
		return m, nil
	}
	return m, nil
}

// Render the stars in a simple fixed grid
func renderStarsGrid(stars []Star) string {
	cols := 20
	rows := 10

	// Create a 2D grid as a map
	pos := make(map[int]string, len(stars))
	for _, s := range stars {
		if s.Visible && s.X >= 0 && s.X < cols && s.Y >= 0 && s.Y < rows {
			key := s.Y*cols + s.X
			pos[key] = s.Symbol
		}
	}

	// Build the grid as a string
	var sb strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if sym, ok := pos[r*cols+c]; ok {
				sb.WriteString(sym)
			} else {
				sb.WriteString("  ") // empty space
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// View function renders the header + stars
func (m model) View() string {
	header := fmt.Sprintf("Theme: %s | Press 1,2,3 to change | q to quit\n\n", themes[m.theme].Name)
	content := renderStarsGrid(m.stars)
	return styleForTheme(themes[m.theme]).Render(header + content)
}

// Main entry point
func main() {
	rand.Seed(time.Now().UnixNano())
	initial := model{
		theme: 0,
		stars: initStars(themes[0]),
	}
	p := tea.NewProgram(initial, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
