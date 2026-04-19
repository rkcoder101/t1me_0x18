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

type CategoryListModel struct {
	client     *api.Client
	categories []api.TaskCategory
	width      int
	height     int
	cursor     int
	loading    bool
	err        error
}

func NewCategoryList(client *api.Client) *CategoryListModel {
	return &CategoryListModel{
		client:     client,
		categories: nil,
		width:      50,
		height:     20,
		cursor:     0,
	}
}

func (m *CategoryListModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *CategoryListModel) fetchCategories() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		categories, err := m.client.GetTaskCategories(ctx)
		if err != nil {
			return CategoryListErrMsg{Err: err}
		}
		return CategoryListLoadedMsg{Categories: categories}
	}
}

type CategoryListLoadedMsg struct {
	Categories []api.TaskCategory
}

type CategoryListErrMsg struct {
	Err error
}

func (m *CategoryListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case CategoryListLoadedMsg:
		m.categories = msg.Categories
		m.loading = false
		if m.cursor >= len(m.categories) {
			m.cursor = 0
		}

	case CategoryListErrMsg:
		m.err = msg.Err
		m.loading = false

	case tea.KeyMsg:
		if m.loading || len(m.categories) == 0 {
			return m, nil
		}

		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *CategoryListModel) View() string {
	if m.err != nil {
		return styles.StyleRed.Render("Error: " + m.err.Error())
	}

	if m.loading {
		return styles.StyleDim.Render("Loading categories...")
	}

	if len(m.categories) == 0 {
		return styles.StyleDim.Render("No categories found")
	}

	var b strings.Builder
	b.WriteString(styles.StyleGreen.Render("Task Categories"))
	b.WriteString("\n")
	b.WriteString(styles.StyleDim.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	for i, cat := range m.categories {
		prefix := "  "
		if i == m.cursor {
			prefix = styles.StyleGreen.Render("› ")
		}

		flex := fmt.Sprintf("flex:%s", cat.SchedulingFlexibility)
		energy := fmt.Sprintf("energy:%s", cat.EnergyRequired)

		row := fmt.Sprintf("%s%s  %s  %s", prefix, cat.Name, flex, energy)

		if i == m.cursor {
			row = styles.SelectedBackground.Render(row)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}

func (m *CategoryListModel) SelectedCategory() *api.TaskCategory {
	if m.cursor >= 0 && m.cursor < len(m.categories) {
		return &m.categories[m.cursor]
	}
	return nil
}

func (m *CategoryListModel) GetCursor() int {
	return m.cursor
}
