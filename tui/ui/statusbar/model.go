package statusbar

import (
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width          int
	DoneCount      int
	RemainingCount int
	OverdueCount   int
}

func New() Model {
	return Model{}
}

func (m Model) View() string {
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.StyleStats(),
		"   ",
		m.StyleHints(),
	)

	return lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		Render(content)
}

func (m Model) StyleStats() string {
	return lipgloss.NewStyle().
		Render(
			m.RenderDone() + " · " + m.RenderRemaining() + " · " + m.RenderOverdue(),
		)
}

func (m Model) RenderDone() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7fff7f"))
	return style.Render("3 done") // placeholder
}

func (m Model) RenderRemaining() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb347"))
	return style.Render("5 remaining") // placeholder
}

func (m Model) RenderOverdue() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	return style.Render("0 overdue") // placeholder
}

func (m Model) StyleHints() string {
	return "[/] command  [?] help"
}

func (m Model) UpdateWidth(width int) Model {
	m.Width = width
	return m
}
