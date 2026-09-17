package core

import "testing"

func testBoard(mines []point) *board {
	b := &board{minesPlaced: true}
	for _, m := range mines {
		b.cells[m.y][m.x].mine = true
	}
	return b
}

func TestAdjacentMines(t *testing.T) {
	b := testBoard([]point{{0, 0}, {2, 0}})

	tests := []struct {
		name string
		x, y int
		want int
	}{
		{"between two mines", 1, 0, 2},
		{"diagonal to two mines", 1, 1, 2},
		{"next to one mine", 0, 1, 1},
		{"far corner", 8, 8, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := b.adjacentMines(tt.x, tt.y); got != tt.want {
				t.Errorf("adjacentMines(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestRevealFlood(t *testing.T) {
	b := testBoard([]point{{0, 0}, {2, 0}})

	if hit := b.reveal(8, 8); hit {
		t.Fatal("reveal on a safe cell reported a mine")
	}

	revealed := 0
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if b.cells[y][x].revealed {
				revealed++
			}
		}
	}

	if want := boardWidth*boardHeight - 3; revealed != want {
		t.Errorf("revealed %d cells, want %d", revealed, want)
	}
	for _, m := range []point{{0, 0}, {2, 0}} {
		if b.cells[m.y][m.x].revealed {
			t.Errorf("mine at %v was revealed by the flood", m)
		}
	}
}

func TestRevealMine(t *testing.T) {
	b := testBoard([]point{{3, 3}})

	if hit := b.reveal(3, 3); !hit {
		t.Error("reveal on a mine cell should report a hit")
	}
}

func TestRevealFlagged(t *testing.T) {
	b := testBoard([]point{{3, 3}})
	b.toggleFlag(3, 3)

	if hit := b.reveal(3, 3); hit {
		t.Error("reveal on a flagged cell should not detonate")
	}
	if b.cells[3][3].revealed {
		t.Error("flagged cell should stay hidden")
	}
}

func TestToggleFlag(t *testing.T) {
	b := testBoard(nil)

	b.toggleFlag(1, 1)
	if !b.cells[1][1].flagged {
		t.Error("expected cell flagged after first toggle")
	}
	b.toggleFlag(1, 1)
	if b.cells[1][1].flagged {
		t.Error("expected flag cleared after second toggle")
	}
	if got := b.remainingMines(); got != mineCount {
		t.Errorf("remainingMines() = %d, want %d", got, mineCount)
	}

	b.toggleFlag(2, 2)
	if got := b.remainingMines(); got != mineCount-1 {
		t.Errorf("remainingMines() = %d, want %d", got, mineCount-1)
	}

	b.cells[4][4].revealed = true
	b.toggleFlag(4, 4)
	if b.cells[4][4].flagged {
		t.Error("revealed cell must not accept a flag")
	}
}

func TestFirstClickSafety(t *testing.T) {
	for trial := 0; trial < 50; trial++ {
		b := &board{}
		if hit := b.reveal(4, 4); hit {
			t.Fatalf("trial %d: first click detonated", trial)
		}
		mines := 0
		for y := 0; y < boardHeight; y++ {
			for x := 0; x < boardWidth; x++ {
				if b.cells[y][x].mine {
					mines++
					if abs(x-4) <= 1 && abs(y-4) <= 1 {
						t.Fatalf("trial %d: mine placed next to first click at (%d, %d)", trial, x, y)
					}
				}
			}
		}
		if mines != mineCount {
			t.Fatalf("trial %d: placed %d mines, want %d", trial, mines, mineCount)
		}
	}
}

func TestCheckWin(t *testing.T) {
	b := testBoard([]point{{0, 0}})
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if x == 0 && y == 0 {
				continue
			}
			b.cells[y][x].revealed = true
		}
	}
	if !b.checkWin() {
		t.Error("expected win when every safe cell is revealed")
	}

	b.cells[8][8].revealed = false
	if b.checkWin() {
		t.Error("expected no win with a safe cell still hidden")
	}
}
