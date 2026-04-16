package config

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

var (
	StyleBase = lipgloss.NewStyle()

	StylePrompt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7fff7f")).
			Render

	StyleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("666666")).
			Render

	StyleAmber = lipgloss.NewStyle().
			Foreground(lipgloss.Color("ffb347")).
			Render

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("2a2a2a"))

	StyleSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("222222")).
			Foreground(lipgloss.Color("7fff7f"))
)

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "scheduler")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}
