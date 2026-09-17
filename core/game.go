package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameState int

const (
	playing gameState = iota
	gameOver
)

type Game struct {
	board  *board
	cursor point
	won    bool
	state  gameState
}

func NewGame() *Game {
	return &Game{board: &board{}, state: playing}
}

func (g *Game) Init() tea.Cmd {
	return nil
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "r":
			if g.state == gameOver {
				return NewGame(), nil
			}
		case "left", "a", "h":
			g.moveCursor(-1, 0)
		case "right", "d", "l":
			g.moveCursor(1, 0)
		case "up", "w", "k":
			g.moveCursor(0, -1)
		case "down", "s", "j":
			g.moveCursor(0, 1)
		case " ", "enter":
			g.revealCursor()
		case "f":
			g.toggleFlagCursor()
		}
	}

	return g, nil
}

func (g *Game) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true).
		MarginBottom(1).
		Render("💣 GoMinesweeper")

	stats := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		MarginBottom(1).
		Render(fmt.Sprintf("Mines: %d", g.board.remainingMines()))

	gameBoard := g.board.render(g.cursor, g.state == gameOver)

	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1).
		Render("Arrows/WASD: Move | Space: Reveal | F: Flag | Q: Quit | R: Restart")

	var statusBar string
	switch g.state {
	case gameOver:
		if g.won {
			statusBar = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD700")).
				Bold(true).
				Render("🎉 YOU WIN | Press R to restart")
		} else {
			statusBar = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true).
				Render("💥 BOOM | Press R to restart")
		}
	default:
		statusBar = ""
	}

	view := fmt.Sprintf("%s\n%s\n\n%s\n%s", title, stats, gameBoard, controls)

	if statusBar != "" {
		view = fmt.Sprintf("%s\n\n%s", view, statusBar)
	}

	return lipgloss.NewStyle().
		MarginLeft(2).
		MarginTop(1).
		Render(view)
}

func (g *Game) moveCursor(dx, dy int) {
	g.cursor.x = clamp(g.cursor.x+dx, 0, boardWidth-1)
	g.cursor.y = clamp(g.cursor.y+dy, 0, boardHeight-1)
}

func clamp(value, lo, hi int) int {
	if value < lo {
		return lo
	}
	if value > hi {
		return hi
	}
	return value
}

func (g *Game) revealCursor() {
	if g.state != playing {
		return
	}
	if g.board.reveal(g.cursor.x, g.cursor.y) {
		g.state = gameOver
		return
	}
	if g.board.checkWin() {
		g.state = gameOver
		g.won = true
	}
}

func (g *Game) toggleFlagCursor() {
	if g.state != playing {
		return
	}
	g.board.toggleFlag(g.cursor.x, g.cursor.y)
}
