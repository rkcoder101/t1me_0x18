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

type RoutineFormModel struct {
	form      *huh.Form
	loading   bool
	err       error
	client    *api.Client
	width     int
	height    int
	editing   bool
	routineID int
	name      string
	weekdays  []string
	startTime string
	duration  string
	isActive  bool
}

func NewRoutineForm(client *api.Client) *RoutineFormModel {
	m := &RoutineFormModel{
		client:    client,
		width:     60,
		height:    20,
		editing:   false,
		routineID: 0,
		name:      "",
		weekdays:  []string{},
		startTime: "09:00",
		duration:  "30",
		isActive:  true,
	}
	m.initForm(nil)
	return m
}

func NewRoutineFormWithRoutine(client *api.Client, routine *api.HardRoutine) *RoutineFormModel {
	weekdays := make([]string, len(routine.Weekdays))
	for i, d := range routine.Weekdays {
		weekdays[i] = string(d)
	}
	m := &RoutineFormModel{
		client:    client,
		width:     60,
		height:    20,
		editing:   true,
		routineID: routine.ID,
		name:      routine.Name,
		weekdays:  weekdays,
		startTime: routine.StartTime,
		duration:  fmt.Sprintf("%d", routine.Duration),
		isActive:  routine.IsActive,
	}
	m.initForm(routine)
	return m
}

func (m *RoutineFormModel) initForm(routine *api.HardRoutine) {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Routine name").
				Value(&m.name),

			huh.NewMultiSelect[string]().
				Title("Weekdays").
				Options(
					huh.NewOption("Monday", "monday"),
					huh.NewOption("Tuesday", "tuesday"),
					huh.NewOption("Wednesday", "wednesday"),
					huh.NewOption("Thursday", "thursday"),
					huh.NewOption("Friday", "friday"),
					huh.NewOption("Saturday", "saturday"),
					huh.NewOption("Sunday", "sunday"),
				).
				Value(&m.weekdays),

			huh.NewInput().
				Title("Start Time").
				Placeholder("HH:MM").
				Value(&m.startTime),

			huh.NewInput().
				Title("Duration (minutes)").
				Placeholder("30").
				Value(&m.duration),

			huh.NewConfirm().
				Title("Active").
				Value(&m.isActive),
		),
	)
}

func (m *RoutineFormModel) Init() tea.Cmd {
	return nil
}

func (m *RoutineFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if !m.loading && m.form != nil {
			form, cmd := m.form.Update(msg)
			m.form = form.(*huh.Form)
			return m, cmd
		}
	}

	return m, nil
}

func (m *RoutineFormModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading...")
	}

	if m.form != nil {
		return m.form.View()
	}

	return styles.StyleDim.Render("Initializing...")
}

type RoutineFormSubmitMsg struct {
	Routine *api.HardRoutine
	Err     error
}

func (m *RoutineFormModel) Submit() tea.Cmd {
	return func() tea.Msg {
		weekdays := make([]api.Weekday, len(m.weekdays))
		for i, d := range m.weekdays {
			weekdays[i] = api.Weekday(d)
		}

		duration := 30
		fmt.Sscanf(m.duration, "%d", &duration)

		routine := &api.HardRoutine{
			Name:      m.name,
			Weekdays:  weekdays,
			StartTime: m.startTime,
			Duration:  duration,
			IsActive:  m.isActive,
		}

		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if m.editing {
			_, err = m.client.UpdateHardRoutine(ctx, m.routineID, routine)
		} else {
			_, err = m.client.CreateHardRoutine(ctx, routine)
		}

		return RoutineFormSubmitMsg{Routine: routine, Err: err}
	}
}

func (m *RoutineFormModel) GetRoutineID() int {
	return m.routineID
}

func (m *RoutineFormModel) IsEditing() bool {
	return m.editing
}
