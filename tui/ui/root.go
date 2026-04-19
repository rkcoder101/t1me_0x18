package ui

import (
	"t1me-tui/api"
	"t1me-tui/ui/dashboard"
	"t1me-tui/ui/palette"
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
	ViewStats
)

type Model struct {
	client      *api.Client
	activeView  ActiveView
	width       int
	height      int
	showPalette bool
	promptInput promptinput.Model
	statusBar   statusbar.Model
	dashboard   dashboard.Model
	palette     palette.Model
}

func New() Model {
	client := api.NewClient("http://localhost:8000")
	return Model{
		client:      client,
		activeView:  ViewDashboard,
		showPalette: false,
		promptInput: promptinput.New(),
		statusBar:   statusbar.New(),
		dashboard:   dashboard.New(client),
		palette:     palette.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.dashboard.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar = m.statusBar.UpdateWidth(msg.Width)
		m.promptInput = m.promptInput.UpdateWidth(msg.Width)
		m.palette = m.palette.UpdateWidth(msg.Width)

	case tea.KeyMsg:
		// Palette overlay - intercept all keys before they reach the dashboard
		if m.showPalette {
			switch msg.String() {
			case "esc":
				m.showPalette = false
				return m, nil
			case "enter":
				if cmd := m.palette.SelectedCommand(); cmd != nil {
					m.showPalette = false
					if cmd.Action == "quit" {
						return m, tea.Quit
					}
					// TODO: Handle other commands (open forms, etc.)
				}
				m.showPalette = false
				return m, nil
			case "ctrl+p":
				m.showPalette = false
				return m, nil
			default:
				m.palette, _ = m.palette.Update(msg)
				return m, nil
			}
		}

		// Handle prompt input focus
		if m.promptInput.IsFocused() {
			switch msg.String() {
			case "esc":
				m.promptInput = m.promptInput.Blur().Clear()
				return m, nil
			case "enter":
				m.promptInput = m.promptInput.Clear().Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.promptInput, cmd = m.promptInput.Update(msg)
			return m, cmd
		}

		// Global keybindings
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+p":
			m.showPalette = true
			return m, nil

		case "ctrl+f":
			var cmd tea.Cmd
			m.promptInput, cmd = m.promptInput.Focus()
			return m, cmd
		}

		// Dashboard handles its own keys
		if m.activeView == ViewDashboard {
			m.dashboard, _ = m.dashboard.Update(msg)
		}
	}

	// Always update dashboard in case of other messages like WindowSize
	if m.activeView == ViewDashboard {
		m.dashboard, _ = m.dashboard.Update(msg)
	}

	return m, nil
}

func (m Model) View() string {
	view := m.renderActiveView()

	// Overlay palette if active - use fixed small box centered
	if m.showPalette {
		paletteView := m.palette.View()
		view = lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, paletteView)
	}

	prompt := m.promptInput.View()
	status := m.statusBar.View()

	promptHeight := lipgloss.Height(prompt)
	statusHeight := lipgloss.Height(status)
	contentHeight := m.height - promptHeight - statusHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width).
		Render(view)

	return content + "\n" + prompt + "\n" + status
}

func (m Model) renderActiveView() string {
	centered := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#999999"))

	switch m.activeView {
	case ViewOnboarding:
		return centered.Render("[Onboarding]")
	case ViewDashboard:
		return m.dashboard.View()
	case ViewTimer:
		return centered.Render("[Timer]")
	case ViewForm:
		return centered.Render("[Form]")
	case ViewStats:
		return centered.Render("[Stats]")
	default:
		return centered.Render("[Unknown]")
	}
}
