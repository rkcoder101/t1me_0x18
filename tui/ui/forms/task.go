package forms

import (
	"context"
	"fmt"
	"time"

	"t1me-tui/api"
	"t1me-tui/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type TaskFormModel struct {
	form          *huh.Form
	loading       bool
	err           error
	client        *api.Client
	categories    []api.TaskCategory
	width         int
	height        int
	editing       bool
	taskID        int
	title         string
	description   string
	categoryID    int
	flexibility   string
	energy        string
	duration      string
	priority      int
	scheduledDate string
}

func NewTaskForm(client *api.Client) *TaskFormModel {
	m := &TaskFormModel{
		client:        client,
		categories:    nil,
		width:         60,
		height:        20,
		editing:       false,
		taskID:        0,
		title:         "",
		description:   "",
		categoryID:    0,
		flexibility:   "M",
		energy:        "M",
		duration:      "30",
		priority:      3,
		scheduledDate: time.Now().Format("2006-01-02"),
	}
	m.initForm(nil)
	return m
}

func NewTaskFormWithTask(client *api.Client, task *api.Task) *TaskFormModel {
	m := &TaskFormModel{
		client:        client,
		categories:    nil,
		width:         60,
		height:        20,
		editing:       true,
		taskID:        task.ID,
		title:         task.Title,
		description:   task.Description,
		categoryID:    task.CategoryID,
		flexibility:   string(task.Flexibility),
		energy:        string(task.EnergyRequired),
		duration:      fmt.Sprintf("%d", task.EstimatedDuration),
		priority:      task.Priority,
		scheduledDate: task.ScheduledDate,
	}
	m.initForm(task)
	return m
}

func (m *TaskFormModel) initForm(task *api.Task) {
	catOptions := make([]huh.Option[int], len(m.categories))
	for i, cat := range m.categories {
		catOptions[i] = huh.NewOption(cat.Name, cat.ID)
	}
	if len(catOptions) == 0 {
		catOptions = []huh.Option[int]{{Key: "No categories", Value: 0}}
	}

	flexOpts := []huh.Option[string]{
		huh.NewOption("Low", "L"),
		huh.NewOption("Medium", "M"),
		huh.NewOption("High", "H"),
	}

	priOpts := []huh.Option[int]{
		huh.NewOption("1 (Highest)", 1),
		huh.NewOption("2", 2),
		huh.NewOption("3 (Default)", 3),
		huh.NewOption("4", 4),
		huh.NewOption("5 (Lowest)", 5),
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Placeholder("Task title").
				Value(&m.title),

			huh.NewInput().
				Title("Description").
				Placeholder("Optional description").
				Value(&m.description),

			huh.NewSelect[int]().
				Title("Category").
				Options(catOptions...).
				Value(&m.categoryID),

			huh.NewSelect[string]().
				Title("Flexibility").
				Options(flexOpts...).
				Value(&m.flexibility),

			huh.NewSelect[string]().
				Title("Energy Required").
				Options(flexOpts...).
				Value(&m.energy),

			huh.NewInput().
				Title("Estimated Duration (minutes)").
				Placeholder("30").
				Value(&m.duration),

			huh.NewSelect[int]().
				Title("Priority").
				Options(priOpts...).
				Value(&m.priority),

			huh.NewInput().
				Title("Scheduled Date").
				Placeholder("YYYY-MM-DD").
				Value(&m.scheduledDate),
		),
	)
}

func (m *TaskFormModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *TaskFormModel) fetchCategories() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		categories, err := m.client.GetTaskCategories(ctx)
		if err != nil {
			return CategoriesErrMsg{Err: err}
		}
		return CategoriesLoadedMsg{Categories: categories}
	}
}

type CategoriesLoadedMsg struct {
	Categories []api.TaskCategory
}

type CategoriesErrMsg struct {
	Err error
}

type TaskFormSubmitMsg struct {
	Task *api.Task
	Err  error
}

func (m *TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case CategoriesLoadedMsg:
		m.categories = msg.Categories
		m.loading = false
		m.initForm(nil)

	case CategoriesErrMsg:
		m.err = msg.Err
		m.loading = false

	case tea.KeyMsg:
		if !m.loading && m.form != nil {
			form, cmd := m.form.Update(msg)
			m.form = form.(*huh.Form)
			return m, cmd
		}
	}

	return m, nil
}

func (m *TaskFormModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading categories...")
	}

	if m.form != nil {
		return m.form.View()
	}

	return styles.StyleDim.Render("Initializing...")
}

func (m *TaskFormModel) Submit() tea.Cmd {
	return func() tea.Msg {
		duration := 30
		fmt.Sscanf(m.duration, "%d", &duration)

		task := &api.Task{
			Title:             m.title,
			Description:       m.description,
			CategoryID:        m.categoryID,
			Flexibility:       api.Flexibility(m.flexibility),
			EnergyRequired:    api.Energy(m.energy),
			EstimatedDuration: duration,
			Priority:          m.priority,
			ScheduledDate:     m.scheduledDate,
			Status:            api.StatusPending,
		}

		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if m.editing {
			_, err = m.client.UpdateTask(ctx, m.taskID, task)
		} else {
			_, err = m.client.CreateTask(ctx, task)
		}

		return TaskFormSubmitMsg{Task: task, Err: err}
	}
}

func (m *TaskFormModel) GetTaskID() int {
	return m.taskID
}

func (m *TaskFormModel) IsEditing() bool {
	return m.editing
}
