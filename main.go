package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"gominesweeper/core"
)

func main() {
	p := tea.NewProgram(core.NewGame(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
