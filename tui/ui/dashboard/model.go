package dashboard

import (
	"context"
	"fmt"
	"strings"
	"time"

	"t1me-tui/api"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardLoadedMsg struct {
	Data *api.DashboardResponse
}

type ErrMsg struct {
	Err error
}

type Model struct {
	client  *api.Client
	data    *api.DashboardResponse
	loading bool
	err     error
	cursor  int
	width   int
	height  int
	spinner spinner.Model
}

func New(client *api.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("7fff7f"))

	return Model{
		client:  client,
		loading: true,
		spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchDashboard())
}

func (m Model) fetchDashboard() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		data, err := m.client.GetDashboardToday(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return DashboardLoadedMsg{Data: data}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case DashboardLoadedMsg:
		m.data = msg.Data
		m.loading = false
		m.cursor = 0 // reset cursor
		return m, nil

	case ErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		totalItems := 0
		if m.data != nil {
			totalItems = len(m.data.Timeline)
		}

		if totalItems == 0 {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.skipNonSelectable(-1)
			}
		case "down", "j":
			if m.cursor < totalItems-1 {
				m.cursor++
				m.skipNonSelectable(1)
			}
		}
	}
	return m, nil
}

func (m *Model) skipNonSelectable(direction int) {
	if m.data == nil {
		return
	}

	timelineLen := len(m.data.Timeline)
	totalItems := timelineLen

	for m.cursor >= 0 && m.cursor < timelineLen {
		item := m.data.Timeline[m.cursor]
		if item.Type == api.TimelineRoutine || item.Type == api.TimelineGap {
			m.cursor += direction
		} else {
			break
		}
	}

	if m.cursor < 0 {
		// Scrolled past the top, find the first selectable or reset to 0
		m.cursor = 0
		for m.cursor < timelineLen {
			item := m.data.Timeline[m.cursor]
			if item.Type != api.TimelineRoutine && item.Type != api.TimelineGap {
				break
			}
			m.cursor++
		}
	} else if m.cursor >= totalItems {
		m.cursor = totalItems - 1
	}
}

func (m Model) UpdateWidth(width int) Model {
	m.width = width
	return m
}

func (m Model) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff7e7e")).Render("Error: " + m.err.Error())
	}

	if m.loading {
		return fmt.Sprintf("%s Loading dashboard...", m.spinner.View())
	}

	if m.data == nil {
		return "No data."
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7fff7f"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	now := time.Now()
	headerText := fmt.Sprintf("%s ─ %s", headerStyle.Render("Today"), dimStyle.Render(now.Format("Monday, Jan 02")))
	timeText := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb347")).Render(now.Format("03:04 PM"))

	// Poor man's align right
	spaces := m.width - lipgloss.Width(headerText) - lipgloss.Width(timeText)
	if spaces < 0 {
		spaces = 0
	}

	b.WriteString(headerText)
	b.WriteString(strings.Repeat(" ", spaces))
	b.WriteString(timeText)
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n\n")

	// Timeline
	for i, item := range m.data.Timeline {
		b.WriteString(m.renderTimelineItem(item, i == m.cursor, now))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderTimelineItem(item api.TimelineItem, selected bool, now time.Time) string {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	if item.Type == api.TimelineGap {
		return dimStyle.Render(fmt.Sprintf("      %s", item.Title))
	}

	timeStr := item.StartTime.Local().Format("15:04")

	isActive := false
	if item.Type == api.TimelineTask {
		endTime := item.StartTime.Add(time.Duration(item.Duration) * time.Minute)
		if now.After(item.StartTime) && now.Before(endTime) {
			isActive = true
		}
	}

	prefix := "  "
	if selected {
		prefix = "│ "
	}
	if isActive {
		prefix = lipgloss.NewStyle().Foreground(lipgloss.Color("#7fff7f")).Render("► ")
	}

	statusIcon := " "
	if item.Status != nil && *item.Status == api.StatusCompleted {
		statusIcon = dimStyle.Render("✓")
		timeStr = dimStyle.Render(timeStr)
	}

	title := item.Title
	if item.Type == api.TimelineRoutine {
		title = dimStyle.Render(title)
	}

	cat := ""
	if item.Type == api.TimelineRoutine {
		cat = dimStyle.Render("[routine]")
	} else if item.CategoryName != nil {
		cat = fmt.Sprintf("[%s]", *item.CategoryName)
	}

	dur := fmt.Sprintf("%dmin", item.Duration)
	if item.Type == api.TimelineRoutine {
		dur = dimStyle.Render(dur)
	}

	// Format: prefix status time title [cat] dur
	return fmt.Sprintf("%s %s  %s  %-25s %-10s %s", prefix, statusIcon, timeStr, title, cat, dur)
}
