package core

import "testing"

func TestMoveCursor(t *testing.T) {
	tests := []struct {
		name   string
		start  point
		dx, dy int
		want   point
	}{
		{"move right", point{4, 4}, 1, 0, point{5, 4}},
		{"move up-left", point{4, 4}, -1, -1, point{3, 3}},
		{"clamp left edge", point{0, 4}, -1, 0, point{0, 4}},
		{"clamp bottom edge", point{4, 8}, 0, 1, point{4, 8}},
		{"clamp corner", point{0, 0}, -5, -5, point{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			g.cursor = tt.start
			g.moveCursor(tt.dx, tt.dy)
			if !g.cursor.equals(tt.want) {
				t.Errorf("cursor = %v, want %v", g.cursor, tt.want)
			}
		})
	}
}

func TestRevealCursorHitsMine(t *testing.T) {
	g := NewGame()
	g.board = testBoard([]point{{0, 0}})
	g.cursor = point{0, 0}

	g.revealCursor()

	if g.state != gameOver {
		t.Errorf("state = %v, want gameOver after hitting a mine", g.state)
	}
	if g.won {
		t.Error("hitting a mine must not set won")
	}
}

func TestRevealCursorWins(t *testing.T) {
	g := NewGame()
	g.board = testBoard([]point{{8, 8}})
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if x == 8 && y == 8 {
				continue
			}
			g.board.cells[y][x].revealed = true
		}
	}
	g.cursor = point{0, 0}

	g.revealCursor()

	if g.state != gameOver || !g.won {
		t.Errorf("state = %v, won = %v; want gameOver and won", g.state, g.won)
	}
}

func TestRevealCursorFlagged(t *testing.T) {
	g := NewGame()
	g.board = testBoard([]point{{0, 0}})
	g.board.toggleFlag(0, 0)
	g.cursor = point{0, 0}

	g.revealCursor()

	if g.state != playing {
		t.Errorf("state = %v, want playing (flagged cell is safe to click)", g.state)
	}
}

func TestInputIgnoredAfterGameOver(t *testing.T) {
	g := NewGame()
	g.board = testBoard([]point{{0, 0}})
	g.cursor = point{0, 0}
	g.revealCursor()

	g.cursor = point{5, 5}
	g.revealCursor()
	g.toggleFlagCursor()

	if g.board.cells[5][5].revealed || g.board.cells[5][5].flagged {
		t.Error("board must not change after game over")
	}
}
