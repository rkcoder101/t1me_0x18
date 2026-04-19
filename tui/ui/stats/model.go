package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"t1me-tui/api"
	"t1me-tui/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type StatsModel struct {
	client  *api.Client
	data    *api.DashboardResponse
	width   int
	height  int
	loading bool
	err     error
}

func NewStats(client *api.Client) *StatsModel {
	return &StatsModel{
		client: client,
		data:   nil,
		width:  50,
		height: 20,
	}
}

func (m *StatsModel) Init() tea.Cmd {
	return m.fetchStats()
}

func (m *StatsModel) fetchStats() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		data, err := m.client.GetDashboardToday(ctx)
		if err != nil {
			return StatsErrMsg{Err: err}
		}
		return StatsLoadedMsg{Data: data}
	}
}

type StatsLoadedMsg struct {
	Data *api.DashboardResponse
}

type StatsErrMsg struct {
	Err error
}

func (m *StatsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case StatsLoadedMsg:
		m.data = msg.Data
		m.loading = false

	case StatsErrMsg:
		m.err = msg.Err
		m.loading = false
	}

	return m, nil
}

func (m *StatsModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading stats...")
	}

	if m.data == nil {
		return styles.StyleDim.Render("No data available")
	}

	var b strings.Builder
	b.WriteString(styles.StyleGreen.Render("Productivity Stats"))
	b.WriteString("\n")
	b.WriteString(styles.StyleDim.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")
	b.WriteString("\n")

	b.WriteString("Today:")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %d tasks", styles.StyleGreen.Render("✓ Done"), m.data.StatsDone))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %d tasks remaining", styles.StyleAmber.Render("○ Remaining"), m.data.StatsRemaining))
	b.WriteString("\n")

	if m.data.StatsOverdue > 0 {
		b.WriteString(fmt.Sprintf("  %s %d overdue", styles.StyleRed.Render("! Overdue"), m.data.StatsOverdue))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.StyleDim.Render("[enter] refresh   [esc] close"))

	return b.String()
}
