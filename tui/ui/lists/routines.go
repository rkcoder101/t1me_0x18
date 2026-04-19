package lists

import (
	"context"
	"fmt"
	"strings"
	"time"

	"t1me-tui/api"
	"t1me-tui/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type RoutineListModel struct {
	client   *api.Client
	routines []api.HardRoutine
	width    int
	height   int
	cursor   int
	loading  bool
	err      error
}

func NewRoutineList(client *api.Client) *RoutineListModel {
	return &RoutineListModel{
		client:   client,
		routines: nil,
		width:    50,
		height:   20,
		cursor:   0,
	}
}

func (m *RoutineListModel) Init() tea.Cmd {
	return m.fetchRoutines()
}

func (m *RoutineListModel) fetchRoutines() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		routines, err := m.client.GetHardRoutines(ctx)
		if err != nil {
			return RoutineListErrMsg{Err: err}
		}
		return RoutineListLoadedMsg{Routines: routines}
	}
}

type RoutineListLoadedMsg struct {
	Routines []api.HardRoutine
}

type RoutineListErrMsg struct {
	Err error
}

func (m *RoutineListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case RoutineListLoadedMsg:
		m.routines = msg.Routines
		m.loading = false
		if m.cursor >= len(m.routines) {
			m.cursor = 0
		}

	case RoutineListErrMsg:
		m.err = msg.Err
		m.loading = false

	case tea.KeyMsg:
		if m.loading || len(m.routines) == 0 {
			return m, nil
		}

		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.routines)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *RoutineListModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading routines...")
	}

	if len(m.routines) == 0 {
		return styles.StyleDim.Render("No routines found")
	}

	var b strings.Builder
	b.WriteString(styles.StyleGreen.Render("Hard Routines"))
	b.WriteString("\n")
	b.WriteString(styles.StyleDim.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	for i, routine := range m.routines {
		prefix := "  "
		if i == m.cursor {
			prefix = styles.StyleGreen.Render("› ")
		}

		active := ""
		if routine.IsActive {
			active = " [active]"
		}

		duration := fmt.Sprintf("%dmin", routine.Duration)
		row := fmt.Sprintf("%s%s%s  %s  %s", prefix, routine.Name, active, routine.StartTime, duration)

		if i == m.cursor {
			row = styles.SelectedBackground.Render(row)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}

func (m *RoutineListModel) SelectedRoutine() *api.HardRoutine {
	if m.cursor >= 0 && m.cursor < len(m.routines) {
		return &m.routines[m.cursor]
	}
	return nil
}

func (m *RoutineListModel) GetCursor() int {
	return m.cursor
}
