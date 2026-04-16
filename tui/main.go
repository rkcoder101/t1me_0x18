package main

import (
	"fmt"
	"t1me-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := ui.New()

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
