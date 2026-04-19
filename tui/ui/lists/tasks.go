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

type TaskListModel struct {
	client  *api.Client
	tasks   []api.Task
	width   int
	height  int
	cursor  int
	loading bool
	err     error
}

func NewTaskList(client *api.Client) *TaskListModel {
	return &TaskListModel{
		client: client,
		tasks:  nil,
		width:  50,
		height: 20,
		cursor: 0,
	}
}

func (m *TaskListModel) Init() tea.Cmd {
	return m.fetchTasks()
}

func (m *TaskListModel) fetchTasks() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tasks, err := m.client.GetTasks(ctx)
		if err != nil {
			return TaskListErrMsg{Err: err}
		}
		return TaskListLoadedMsg{Tasks: tasks}
	}
}

type TaskListLoadedMsg struct {
	Tasks []api.Task
}

type TaskListErrMsg struct {
	Err error
}

func (m *TaskListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case TaskListLoadedMsg:
		m.tasks = msg.Tasks
		m.loading = false
		if m.cursor >= len(m.tasks) {
			m.cursor = 0
		}

	case TaskListErrMsg:
		m.err = msg.Err
		m.loading = false

	case tea.KeyMsg:
		if m.loading || len(m.tasks) == 0 {
			return m, nil
		}

		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "enter":
			// Selection handled by parent
		}
	}

	return m, nil
}

func (m *TaskListModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading tasks...")
	}

	if len(m.tasks) == 0 {
		return styles.StyleDim.Render("No tasks found")
	}

	var b strings.Builder
	b.WriteString(styles.StyleGreen.Render("Tasks"))
	b.WriteString("\n")
	b.WriteString(styles.StyleDim.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	for i, task := range m.tasks {
		prefix := "  "
		if i == m.cursor {
			prefix = styles.StyleGreen.Render("› ")
		}

		status := ""
		if task.Status != "" {
			status = fmt.Sprintf(" [%s]", task.Status)
		}

		duration := fmt.Sprintf("%dmin", task.EstimatedDuration)
		row := fmt.Sprintf("%s%s%s  %s", prefix, task.Title, status, duration)

		if i == m.cursor {
			row = styles.SelectedBackground.Render(row)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}

func (m *TaskListModel) SelectedTask() *api.Task {
	if m.cursor >= 0 && m.cursor < len(m.tasks) {
		return &m.tasks[m.cursor]
	}
	return nil
}

func (m *TaskListModel) GetCursor() int {
	return m.cursor
}

func (m *TaskListModel) SetCursor(c int) {
	m.cursor = c
}
