package dialog

import (
	"context"
	"fmt"
	"time"

	"t1me-tui/api"
	"t1me-tui/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DialogModel struct {
	title   string
	message string
	width   int
	height  int
	cursor  int
	confirm bool
}

func NewConfirmDialog(title, message string) *DialogModel {
	return &DialogModel{
		title:   title,
		message: message,
		width:   40,
		height:  10,
		cursor:  0,
	}
}

func (m *DialogModel) Init() tea.Cmd {
	return nil
}

func (m *DialogModel) GetConfirm() bool {
	return m.confirm
}

func (m *DialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 4
		m.height = msg.Height - 6

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right":
			m.cursor = 1 - m.cursor
		case "enter":
			m.confirm = m.cursor == 0
		}
	}

	return m, nil
}

func (m *DialogModel) View() string {
	confirmLabel := "[ Confirm ]"
	cancelLabel := "[ Cancel ]"

	if m.cursor == 0 {
		confirmLabel = styles.SelectedBackground.Render(confirmLabel)
	} else {
		cancelLabel = styles.SelectedBackground.Render(cancelLabel)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Green).
		Width(m.width).
		Padding(1, 2)

	content := fmt.Sprintf("%s\n\n%s\n\n%s  %s",
		styles.StyleGreen.Render(m.title),
		m.message,
		confirmLabel,
		cancelLabel)

	return box.Render(content)
}

type DeleteDialogModel struct {
	itemType string
	itemName string
	itemID   int
	width    int
	height   int
	cursor   int
	confirm  bool
	client   *api.Client
}

func NewDeleteDialog(itemType string, itemName string, itemID int, client *api.Client) *DeleteDialogModel {
	return &DeleteDialogModel{
		itemType: itemType,
		itemName: itemName,
		itemID:   itemID,
		width:    40,
		height:   10,
		cursor:   0,
		client:   client,
	}
}

func (m *DeleteDialogModel) Init() tea.Cmd {
	return nil
}

func (m *DeleteDialogModel) GetConfirm() bool {
	return m.confirm
}

func (m *DeleteDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 4
		m.height = msg.Height - 6

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right":
			m.cursor = 1 - m.cursor
		case "enter":
			m.confirm = m.cursor == 0
		}
	}

	return m, nil
}

func (m *DeleteDialogModel) View() string {
	confirmLabel := "[ Delete ]"
	cancelLabel := "[ Cancel ]"

	if m.cursor == 0 {
		confirmLabel = styles.StyleRed.Render("[ Delete ]")
	} else {
		cancelLabel = styles.SelectedBackground.Render(cancelLabel)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Red).
		Width(m.width).
		Padding(1, 2)

	message := fmt.Sprintf("Delete %s '%s'?", m.itemType, m.itemName)

	content := fmt.Sprintf("%s\n\nThis action cannot be undone.\n\n%s  %s",
		styles.StyleRed.Render("Delete "+m.itemType),
		message,
		confirmLabel,
		cancelLabel)

	return box.Render(content)
}

func (m *DeleteDialogModel) DeleteItem() error {
	if !m.confirm {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch m.itemType {
	case "task":
		return m.client.DeleteTask(ctx, m.itemID)
	case "routine":
		return m.client.DeleteHardRoutine(ctx, m.itemID)
	case "category":
		return m.client.DeleteTaskCategory(ctx, m.itemID)
	}
	return nil
}
