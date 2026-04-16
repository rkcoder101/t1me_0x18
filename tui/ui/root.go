package ui

import (
	"t1me-tui/ui/promptinput"
	"t1me-tui/ui/statusbar"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveView int

const (
	ViewOnboarding ActiveView = iota
	ViewDashboard
	ViewTimer
	ViewForm
)

type Model struct {
	activeView  ActiveView
	width       int
	height      int
	promptInput promptinput.Model
	statusBar   statusbar.Model
	showPalette bool
}

func New() Model {
	return Model{
		activeView:  ViewDashboard,
		promptInput: promptinput.New(),
		statusBar:   statusbar.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar = m.statusBar.UpdateWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		// Layer 1: Global Overrides (works from anywhere)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			// TODO: Force refresh from API
			return m, nil
		}

		// Layer 2: Active Overlay (Palette)
		if m.showPalette {
			switch msg.String() {
			case "esc":
				m.showPalette = false
				return m, nil
			case "ctrl+p":
				m.showPalette = false
				return m, nil
			}
			// TODO: Pass to palette model for navigation/selection
			return m, nil
		}

		// Layer 3: Focused Element (AI Prompt Input)
		if m.promptInput.IsFocused() {
			switch msg.String() {
			case "esc":
				m.promptInput = m.promptInput.Blur().Clear()
				return m, nil
			case "enter":
				value := m.promptInput.Value()
				if value != "" {
					// TODO: Fire AI prompt submission via API
					m.promptInput = m.promptInput.Clear().Blur()
				}
				return m, nil
			case "ctrl+p":
				m.promptInput = m.promptInput.Blur()
				m.showPalette = true
				return m, nil
			}
			var cmd tea.Cmd
			m.promptInput, cmd = m.promptInput.Update(msg)
			return m, cmd
		}

		// Layer 4: Active View (Dashboard, Timer, etc.)
		switch msg.String() {
		case "ctrl+p":
			m.showPalette = true
			return m, nil
		case "ctrl+s":
			var cmd tea.Cmd
			m.promptInput, cmd = m.promptInput.Focus()
			return m, cmd
		case "1":
			m.activeView = ViewOnboarding
		case "2":
			m.activeView = ViewDashboard
		case "3":
			m.activeView = ViewTimer
		case "4":
			m.activeView = ViewForm
		}
	}

	return m, nil
}

func (m Model) View() string {
	view := m.renderActiveView()
	prompt := m.promptInput.View()
	status := m.statusBar.View()

	contentHeight := m.height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width).
		Render(view)

	// If palette is shown, overlay it
	if m.showPalette {
		paletteView := m.renderPalette()
		content = lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, paletteView)
	}

	return content + "\n" + prompt + "\n" + status
}

func (m Model) renderActiveView() string {
	centered := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#999999"))

	switch m.activeView {
	case ViewOnboarding:
		return centered.Render("[Onboarding view placeholder]\n\nPress 2 for Dashboard")
	case ViewDashboard:
		return centered.Render("[Dashboard view placeholder]\n\nPress 1 for Onboarding\nPress 3 for Timer\nPress 4 for Form")
	case ViewTimer:
		return centered.Render("[Timer view placeholder]\n\nPress 2 for Dashboard")
	case ViewForm:
		return centered.Render("[Form view placeholder]\n\nPress 2 for Dashboard")
	default:
		return centered.Render("[Unknown view]")
	}
}

func (m Model) renderPalette() string {
	// Placeholder for command palette
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("2a2a2a")).
		Padding(1, 2).
		Width(40)

	return box.Render("[Command Palette - press esc to close]")
}
