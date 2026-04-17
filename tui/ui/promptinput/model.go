package promptinput

import (
	"github.com/charmbracelet/bubbles/textarea"
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
	input textarea.Model
}

func New() Model {
	ta := textarea.New()
	ta.Placeholder = "ask AI to schedule something..."
	ta.FocusedStyle.Placeholder = placeholderStyle
	ta.Prompt = "› "
	ta.FocusedStyle.Prompt = promptStyle
	ta.ShowLineNumbers=false
	ta.CharLimit = 500
	ta.SetHeight(1)
	return Model{input: ta}
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

func (m Model) UpdateWidth(width int) Model {
	m.input.SetWidth(width)
	return m
}

func (m Model) Value() string { return m.input.Value() }

func (m Model) View() string {
	return m.input.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	textWidth := m.input.Width() - 2 
    if textWidth <= 0 { textWidth = 1 }
	contentLen := len(m.input.Value())
	lines := (contentLen + textWidth - 1) / textWidth
    if lines < 1 { lines = 1 }
    if lines > 3 { lines = 3 } // Cap it at 3 lines for now
    m.input.SetHeight(lines)
	return m, cmd
}
