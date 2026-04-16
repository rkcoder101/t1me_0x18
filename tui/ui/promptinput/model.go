package promptinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7fff7f"))
	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("666666"))
)

type Model struct {
	input textinput.Model
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "ask AI to schedule something..."
	ti.PlaceholderStyle = placeholderStyle
	ti.Prompt = "› "
	ti.PromptStyle = promptStyle
	ti.CharLimit = 500
	return Model{input: ti}
}

func (m Model) IsFocused() bool { return m.input.Focused() }

func (m Model) Focus() (Model, tea.Cmd) {
	cmd := m.input.Focus()
	return m, cmd
}

func (m Model) Blur() Model {
	m.input.Blur()
	return m
}

func (m Model) Clear() Model {
	m.input.Reset()
	return m
}

func (m Model) Value() string { return m.input.Value() }

func (m Model) View() string {
	return m.input.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
