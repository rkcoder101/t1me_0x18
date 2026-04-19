package palette

import (
	"strings"

	"t1me-tui/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Command struct {
	Label  string
	Action string
}

var DefaultCommands = []Command{
	{Label: "Add Task", Action: "add_task"},
	{Label: "Add Hard Routine", Action: "add_routine"},
	{Label: "Add Task Category", Action: "add_category"},
	{Label: "Edit Task", Action: "edit_task"},
	{Label: "Edit Hard Routine", Action: "edit_routine"},
	{Label: "Edit Task Category", Action: "edit_category"},
	{Label: "Delete Task", Action: "delete_task"},
	{Label: "Delete Hard Routine", Action: "delete_routine"},
	{Label: "Delete Task Category", Action: "delete_category"},
	{Label: "View All Tasks", Action: "view_tasks"},
	{Label: "View Hard Routines", Action: "view_routines"},
	{Label: "View Task Categories", Action: "view_categories"},
	{Label: "View Stats", Action: "view_stats"},
	{Label: "Schedule Today", Action: "schedule_today"},
	{Label: "Shift Tasks", Action: "shift_tasks"},
	{Label: "Settings", Action: "settings"},
	{Label: "Quit", Action: "quit"},
}

const (
	inputLines  = 1
	dividerLine = 1
	hintLines   = 1
	borderLines = 2
)

type Model struct {
	input        textinput.Model
	commands     []Command
	filtered     []Command
	cursor       int
	scrollOffset int
	termWidth    int
	termHeight   int
	Quit         bool
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "filter..."
	input.Prompt = ""
	input.TextStyle = lipgloss.Style{}
	input.Focus()

	return Model{
		input:        input,
		commands:     DefaultCommands,
		filtered:     DefaultCommands,
		cursor:       0,
		scrollOffset: 0,
		termWidth:    80,
		termHeight:   24,
	}
}

func (m Model) UpdateWidth(w int) Model {
	m.termWidth = w
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+p", "esc":
			m.Quit = true
			return m, nil

		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			} else {
				m.cursor = len(m.filtered) - 1
				m.ensureVisible()
			}

		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.ensureVisible()
			} else {
				m.cursor = 0
				m.scrollOffset = 0
			}

		case "enter":
			m.Quit = true

		default:
			m.input, _ = m.input.Update(msg)
			m.filterCommands()
		}
	}

	return m, nil
}

func (m *Model) ensureVisible() {
	maxVisible := m.maxVisible()

	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+maxVisible {
		m.scrollOffset = m.cursor - maxVisible + 1
	}
}

func (m *Model) filterCommands() {
	query := m.input.Value()
	if query == "" {
		m.filtered = m.commands
	} else {
		m.filtered = nil
		for _, cmd := range m.commands {
			if strings.Contains(strings.ToLower(cmd.Label), strings.ToLower(query)) {
				m.filtered = append(m.filtered, cmd)
			}
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = 0
	}
	m.scrollOffset = 0
}

func (m Model) boxWidth() int {
	w := m.termWidth * 3 / 10
	if w < 20 {
		w = 20
	}
	if w > 30 {
		w = 30
	}
	return w
}

func (m Model) maxVisible() int {
	available := m.termHeight - inputLines - dividerLine - hintLines - borderLines - 4
	if available < 4 {
		available = 4
	}
	if available > 10 {
		available = 10
	}
	return available
}

func (m Model) View() string {
	boxW := m.boxWidth()
	maxVis := m.maxVisible()
	totalCmds := len(m.filtered)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Green).
		Width(boxW+2).
		Margin(1, 0)

	var b strings.Builder

	// Filter input
	inputView := strings.TrimSpace(m.input.View())
	b.WriteString(inputView)
	b.WriteString("\n")

	// Divider - use Unicode for nicer look
	sep := "─"
	if m.scrollOffset > 0 {
		sep = "┄"
	}
	b.WriteString(strings.Repeat(sep, boxW))
	b.WriteString("\n")

	// Visible commands
	visible := maxVis
	if totalCmds < visible {
		visible = totalCmds
	}

	for i := 0; i < visible; i++ {
		idx := m.scrollOffset + i
		if idx >= totalCmds {
			break
		}

		row := m.filtered[idx].Label
		padding := boxW - len(row)
		if padding > 0 {
			row += strings.Repeat(" ", padding)
		}

		if idx == m.cursor {
			row = styles.PaletteSelected.Render(row)
		} else {
			row = styles.PaletteNormal.Render(" " + row)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}

	// Empty rows to maintain box size
	for i := visible; i < maxVis; i++ {
		b.WriteString(strings.Repeat(" ", boxW))
		b.WriteString("\n")
	}

	// Scroll indicator / hint
	hint := ""
	if totalCmds > maxVis {
		if m.scrollOffset > 0 && m.scrollOffset+maxVis < totalCmds {
			hint = " (" + styles.StyleGreen.Render("↑↓") + ")"
		} else if m.scrollOffset > 0 {
			hint = " (" + styles.StyleGreen.Render("↑") + ")"
		} else if m.scrollOffset+maxVis < totalCmds {
			hint = " (" + styles.StyleGreen.Render("↓") + ")"
		}
	}
	hint += " [enter] [esc]"
	b.WriteString(hint)

	return box.Render(b.String())
}

func (m Model) SelectedCommand() *Command {
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		return &m.filtered[m.cursor]
	}
	return nil
}

func (m Model) Close() {
	m.Quit = true
}
