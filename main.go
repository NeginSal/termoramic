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

type tickMsg time.Time
type Star struct {
	X       int
	Y       int
	Visible bool
}

const StarCount = 40

type model struct {
	stars  []Star
	theme  int
	width  int
	height int
}

var style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#874BFD")).
	Bold(true).
	Padding(1, 2)

func initStars(w, h int) []Star {
	cols := 40
	rows := 10

	stars := make([]Star, StarCount)
	for i := 0; i < StarCount; i++ {
		stars[i] = Star{
			X:       rand.Intn(cols),
			Y:       rand.Intn(rows),
			Visible: rand.Intn(2) == 0,
		}
	}
	return stars
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		for i := range m.stars {
			if rand.Intn(100) < 30 {
				m.stars[i].Visible = !m.stars[i].Visible
			}
		}
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func renderStarsGrid(stars []Star) string {
	cols := 40
	rows := 10

	pos := make(map[int]bool, len(stars))
	for _, s := range stars {
		if s.Visible {
			key := s.Y*cols + s.X
			pos[key] = true
		}
	}

	var sb strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pos[r*cols+c] {
				sb.WriteString("✨")
			} else {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m model) View() string {
	header := "Press q to quit\n\n"
	content := renderStarsGrid(m.stars)
	return style.Render(header + content)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	initial := model{
		stars: initStars(80, 24),
		theme: 0,
	}
	p := tea.NewProgram(initial, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}