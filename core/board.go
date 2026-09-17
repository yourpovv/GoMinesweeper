package core

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	boardWidth  = 9
	boardHeight = 9
	mineCount   = 10
)

type cell struct {
	mine     bool
	revealed bool
	flagged  bool
}

type board struct {
	cells       [boardHeight][boardWidth]cell
	minesPlaced bool
}

func (b *board) placeMines(safe point) {
	var candidates []point
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if abs(x-safe.x) <= 1 && abs(y-safe.y) <= 1 {
				continue
			}
			candidates = append(candidates, point{x, y})
		}
	}
	for _, i := range rand.Perm(len(candidates))[:mineCount] {
		p := candidates[i]
		b.cells[p.y][p.x].mine = true
	}
	b.minesPlaced = true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (b *board) inBounds(x, y int) bool {
	return x >= 0 && x < boardWidth && y >= 0 && y < boardHeight
}

func (b *board) neighbors(x, y int) []point {
	var out []point
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if p := (point{x + dx, y + dy}); b.inBounds(p.x, p.y) {
				out = append(out, p)
			}
		}
	}
	return out
}

func (b *board) adjacentMines(x, y int) int {
	count := 0
	for _, n := range b.neighbors(x, y) {
		if b.cells[n.y][n.x].mine {
			count++
		}
	}
	return count
}

func (b *board) reveal(x, y int) bool {
	cell := &b.cells[y][x]
	if cell.revealed || cell.flagged {
		return false
	}
	if !b.minesPlaced {
		b.placeMines(point{x, y})
	}
	cell.revealed = true
	if cell.mine {
		return true
	}
	if b.adjacentMines(x, y) > 0 {
		return false
	}
	stack := b.neighbors(x, y)
	for len(stack) > 0 {
		next := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		neighbor := &b.cells[next.y][next.x]
		if neighbor.mine || neighbor.revealed || neighbor.flagged {
			continue
		}
		neighbor.revealed = true
		if b.adjacentMines(next.x, next.y) == 0 {
			stack = append(stack, b.neighbors(next.x, next.y)...)
		}
	}
	return false
}

func (b *board) toggleFlag(x, y int) {
	cell := &b.cells[y][x]
	if cell.revealed {
		return
	}
	cell.flagged = !cell.flagged
}

func (b *board) checkWin() bool {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cell := b.cells[y][x]
			if !cell.mine && !cell.revealed {
				return false
			}
		}
	}
	return true
}

func (b *board) remainingMines() int {
	flagged := 0
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if b.cells[y][x].flagged {
				flagged++
			}
		}
	}
	return mineCount - flagged
}

var numberColors = map[int]string{
	1: "#5B8DEF",
	2: "#4CAF50",
	3: "#F44336",
	4: "#9C27B0",
	5: "#FF9800",
	6: "#00BCD4",
	7: "#E91E63",
	8: "#9E9E9E",
}

func (b *board) render(cursor point, revealAll bool) string {
	var sb strings.Builder

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	sb.WriteString(borderStyle.Render("╔"+strings.Repeat("═", boardWidth*3)+"╗") + "\n")

	for y := 0; y < boardHeight; y++ {
		sb.WriteString(borderStyle.Render("║"))
		for x := 0; x < boardWidth; x++ {
			sb.WriteString(b.cellView(point{x, y}, cursor, revealAll))
		}
		sb.WriteString(borderStyle.Render("║") + "\n")
	}

	sb.WriteString(borderStyle.Render("╚"+strings.Repeat("═", boardWidth*3)+"╝") + "\n")

	return sb.String()
}

func (b *board) cellView(p point, cursor point, revealAll bool) string {
	cell := b.cells[p.y][p.x]

	content, color := "#", "#666666"
	switch {
	case cell.flagged:
		content, color = "F", "#FFFF00"
	case !cell.revealed && !revealAll:
		content, color = "#", "#666666"
	case cell.mine:
		content, color = "*", "#FF0000"
	default:
		count := b.adjacentMines(p.x, p.y)
		if count == 0 {
			content, color = " ", "#666666"
		} else {
			content, color = string(rune('0'+count)), numberColors[count]
		}
	}

	style := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color(color)).
		Bold(true)
	if cursor.equals(p) {
		style = style.Reverse(true)
	}
	return style.Render(content)
}
