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

// tickMsg is emitted by tea.Tick to drive animations.
type tickMsg time.Time

// Grid size kept simple on purpose (no window-size handling for now).
const (
	GridCols = 20
	GridRows = 10
)

// Theme IDs for clarity (must match the themes slice order).
const (
	ThemeStarry = iota
	ThemeFlowers
	ThemeOcean
)

// Star is a renderable item on the grid.
type Star struct {
	X       int    // column (0..GridCols-1)
	Y       int    // row (0..GridRows-1)
	Visible bool   // whether to draw it or not
	Symbol  string // emoji/symbol to print (✨, 🌸, 🌊, ...)
	DX      int    // horizontal velocity (used by Flowers/Ocean)
}

// Number of items to render.
const StarCount = 40

// Theme defines colors and available symbols for a theme.
type Theme struct {
	Name       string
	Background string
	Foreground string
	Symbols    []string
}

// Themes list (order must match Theme* consts).
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

// model holds the program state.
type model struct {
	stars []Star
	theme int // ThemeStarry / ThemeFlowers / ThemeOcean
}

// styleForTheme builds a lipgloss style for the current theme.
func styleForTheme(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Foreground)).
		Background(lipgloss.Color(t.Background)).
		Padding(1, 2)
}

// initStars creates initial items based on the selected theme.
// Simplicity: we use a fixed GridCols x GridRows grid.
func initStars(themeID int) []Star {
	t := themes[themeID]
	stars := make([]Star, StarCount)
	for i := 0; i < StarCount; i++ {
		dx := 0
		switch themeID {
		case ThemeFlowers:
			// Flowers sway left/right: start with -1 or +1 randomly.
			if rand.Intn(2) == 0 {
				dx = -1
			} else {
				dx = 1
			}
		case ThemeOcean:
			// Ocean flows to the right.
			dx = 1
		}

		stars[i] = Star{
			X:       rand.Intn(GridCols),
			Y:       rand.Intn(GridRows),
			Visible: rand.Intn(2) == 0, // ~50% visible initially
			Symbol:  t.Symbols[rand.Intn(len(t.Symbols))],
			DX:      dx,
		}
	}
	return stars
}

// Init schedules the first tick to start animations.
func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles ticks and keypresses, animating per theme.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		// Per-theme animation behavior
		switch m.theme {

		case ThemeStarry:
			// Blink randomly: toggle visibility with small probability.
			for i := range m.stars {
				if rand.Intn(100) < 30 { // 30% chance to flip
					m.stars[i].Visible = !m.stars[i].Visible
				}
			}

		case ThemeFlowers:
			// Flowers sway: move horizontally and bounce at edges.
			for i := range m.stars {
				// Optional subtle blinking for life-like feel
				if rand.Intn(100) < 10 {
					m.stars[i].Visible = !m.stars[i].Visible
				}
				m.stars[i].X += m.stars[i].DX
				// Bounce when hitting edges
				if m.stars[i].X < 0 {
					m.stars[i].X = 0
					m.stars[i].DX = 1
				}
				if m.stars[i].X >= GridCols {
					m.stars[i].X = GridCols - 1
					m.stars[i].DX = -1
				}
			}

		case ThemeOcean:
			// Ocean flows: shift to the right with wrap-around.
			for i := range m.stars {
				// Light shimmer
				if rand.Intn(100) < 10 {
					m.stars[i].Visible = !m.stars[i].Visible
				}
				m.stars[i].X = (m.stars[i].X + 1) % GridCols
			}
		}

		// Schedule the next tick
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.theme = ThemeStarry
			m.stars = initStars(m.theme)
		case "2":
			m.theme = ThemeFlowers
			m.stars = initStars(m.theme)
		case "3":
			m.theme = ThemeOcean
			m.stars = initStars(m.theme)
		}
		return m, nil
	}
	return m, nil
}

// renderStarsGrid draws the current items on a fixed-size grid.
func renderStarsGrid(stars []Star) string {
	cols := GridCols
	rows := GridRows

	// Map occupied cells to the symbol to render.
	pos := make(map[int]string, len(stars))
	for _, s := range stars {
		if s.Visible && s.X >= 0 && s.X < cols && s.Y >= 0 && s.Y < rows {
			key := s.Y*cols + s.X
			pos[key] = s.Symbol
		}
	}

	// Build the grid line by line.
	var sb strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if sym, ok := pos[r*cols+c]; ok {
				sb.WriteString(sym)
			} else {
				sb.WriteString("  ") // empty spaces keep layout aligned
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// View renders the header and the themed content box.
func (m model) View() string {
	header := fmt.Sprintf("Theme: %s | 1:Starry  2:Flowers  3:Ocean | q:quit\n\n", themes[m.theme].Name)
	content := renderStarsGrid(m.stars)
	return styleForTheme(themes[m.theme]).Render(header + content)
}

// main sets up initial state and starts the Bubble Tea program.
func main() {
	rand.Seed(time.Now().UnixNano())

	initial := model{
		theme: ThemeStarry,
		stars: initStars(ThemeStarry),
	}

	p := tea.NewProgram(initial, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
