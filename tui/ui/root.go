package ui

import (
	"t1me-tui/api"
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
	activeView ActiveView
	width      int
	height     int
	client     *api.Client
	prompt     promptinput.Model
	statusBar  statusbar.Model
}

func New(client *api.Client) Model {
	return Model{
		activeView: ViewDashboard,
		client:     client,
		prompt:     promptinput.New(),
		statusBar:  statusbar.New(),
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
		m.prompt = m.prompt.UpdateWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activeView = ViewOnboarding
		case "2":
			m.activeView = ViewDashboard
		case "3":
			m.activeView = ViewTimer
		case "4":
			m.activeView = ViewForm
		case "ctrl+space":
			// TODO: Open command palette
		}

		m.prompt, _ = m.prompt.Update(msg)
	}
	return m, nil
}

func (m Model) View() string {
	prompt := m.prompt.View()
	view := m.renderActiveView()
	statusBar := m.statusBar.View()

	mainContent := lipgloss.NewStyle().
		Height(m.height - 2). // -2 for prompt input and status bar
		Render(view)

	return mainContent + "\n" + prompt + "\n" + statusBar
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
