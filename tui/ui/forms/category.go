package forms

import (
	"context"
	"time"

	"t1me-tui/api"
	"t1me-tui/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CategoryFormModel struct {
	form        *huh.Form
	loading     bool
	err         error
	client      *api.Client
	width       int
	height      int
	editing     bool
	categoryID  int
	name        string
	flexibility string
	energy      string
	needsFocus  bool
}

func NewCategoryForm(client *api.Client) *CategoryFormModel {
	m := &CategoryFormModel{
		client:      client,
		width:       60,
		height:      20,
		editing:     false,
		categoryID:  0,
		name:        "",
		flexibility: "M",
		energy:      "M",
		needsFocus:  false,
	}
	m.initForm(nil)
	return m
}

func NewCategoryFormWithCategory(client *api.Client, category *api.TaskCategory) *CategoryFormModel {
	m := &CategoryFormModel{
		client:      client,
		width:       60,
		height:      20,
		editing:     true,
		categoryID:  category.ID,
		name:        category.Name,
		flexibility: string(category.SchedulingFlexibility),
		energy:      string(category.EnergyRequired),
		needsFocus:  category.NeedsFocusBlock,
	}
	m.initForm(category)
	return m
}

func (m *CategoryFormModel) initForm(category *api.TaskCategory) {
	flexOpts := []huh.Option[string]{
		huh.NewOption("Low", "L"),
		huh.NewOption("Medium", "M"),
		huh.NewOption("High", "H"),
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Category name").
				Value(&m.name),

			huh.NewSelect[string]().
				Title("Scheduling Flexibility").
				Options(flexOpts...).
				Value(&m.flexibility),

			huh.NewSelect[string]().
				Title("Energy Required").
				Options(flexOpts...).
				Value(&m.energy),

			huh.NewConfirm().
				Title("Needs Focus Block").
				Value(&m.needsFocus),
		),
	)
}

func (m *CategoryFormModel) Init() tea.Cmd {
	return nil
}

func (m *CategoryFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *CategoryFormModel) View() string {
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

type CategoryFormSubmitMsg struct {
	Category *api.TaskCategory
	Err      error
}

func (m *CategoryFormModel) Submit() tea.Cmd {
	return func() tea.Msg {
		category := &api.TaskCategory{
			Name:                  m.name,
			SchedulingFlexibility: api.Flexibility(m.flexibility),
			EnergyRequired:        api.Energy(m.energy),
			NeedsFocusBlock:       m.needsFocus,
		}

		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if m.editing {
			_, err = m.client.UpdateTaskCategory(ctx, m.categoryID, category)
		} else {
			_, err = m.client.CreateTaskCategory(ctx, category)
		}

		return CategoryFormSubmitMsg{Category: category, Err: err}
	}
}

func (m *CategoryFormModel) GetCategoryID() int {
	return m.categoryID
}

func (m *CategoryFormModel) IsEditing() bool {
	return m.editing
}
