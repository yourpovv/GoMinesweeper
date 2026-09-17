<div align="center">
  
# goMinesweeper

**Terminal minesweeper game built with Go, Bubbletea and Lipgloss.**

[![GitHub](https://img.shields.io/github/stars/yourpovv/GoMinesweeper?style=social)](https://github.com/yourpovv/GoMinesweeper)

https://github.com/user-attachments/assets/2b8f58d8-91b9-4c55-8c07-4ba76ff7d136

</div>

## Requirements

- Go 1.21+

## Running it

```bash
go run .
```

or build it first:

```bash
go build -o goMinesweeper.exe
.\goMinesweeper.exe
```

## Controls

- Arrow keys or WASD to move the cursor
- Space or Enter to reveal a cell
- F to flag a suspected mine
- Q to quit
- R to restart after game over

## How it works

Clear a 9x9 board hiding 10 mines. Numbers count adjacent mines, empty cells flood open. Your first click is always safe. Flag suspects with F, reveal every safe cell to win. Hit a mine and it's over.
