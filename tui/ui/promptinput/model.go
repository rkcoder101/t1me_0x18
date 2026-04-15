package promptinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Input     textinput.Model
	Width     int
	Submitted bool
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "type a task for AI to schedule..."
	ti.Prompt = "› "
	ti.Focus()

	return Model{
		Input: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.Input, _ = m.Input.Update(msg)
	return m, nil
}

func (m Model) View() string {
	input := m.Input.View()

	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffb347")).
		Render("[enter] send")

	inputWidth := lipgloss.Width(input)
	hintWidth := lipgloss.Width(hint)
	gap := m.Width - inputWidth - hintWidth
	if gap < 2 {
		gap = 2
	}

	return input + lipgloss.NewStyle().PaddingLeft(gap).Render(hint)
}

func (m Model) GetValue() string {
	return m.Input.Value()
}

func (m Model) Clear() Model {
	m.Input.SetValue("")
	m.Submitted = false
	return m
}

func (m Model) UpdateWidth(width int) Model {
	m.Width = width
	return m
}
