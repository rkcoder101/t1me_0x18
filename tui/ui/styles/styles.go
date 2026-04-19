package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Green       = lipgloss.Color("#7fff7f")
	Amber       = lipgloss.Color("#ffb347")
	Blue        = lipgloss.Color("#7eb8ff")
	Red         = lipgloss.Color("#ff7e7e")
	Text        = lipgloss.Color("#e8e8e8")
	Dim         = lipgloss.Color("#666666")
	Mid         = lipgloss.Color("#999999")
	Bg          = lipgloss.Color("#0d0d0d")
	BorderColor = lipgloss.Color("#2a2a2a")

	// Base Styles
	StyleGreen = lipgloss.NewStyle().Foreground(Green)
	StyleAmber = lipgloss.NewStyle().Foreground(Amber)
	StyleBlue  = lipgloss.NewStyle().Foreground(Blue)
	StyleRed   = lipgloss.NewStyle().Foreground(Red)
	StyleText  = lipgloss.NewStyle().Foreground(Text)
	StyleDim   = lipgloss.NewStyle().Foreground(Dim)
	StyleMid   = lipgloss.NewStyle().Foreground(Mid)
	StyleBg    = lipgloss.NewStyle().Background(Bg)

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor)

	StyleSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a1a")).
			Foreground(Blue).
			Bold(true)

	SelectedBackground = lipgloss.NewStyle().
				Background(lipgloss.Color("#222222")).
				Foreground(Green)

	PaletteSelected = lipgloss.NewStyle().
			Background(Green).
			Foreground(lipgloss.Color("#0d0d0d")).
			Bold(true)

	PaletteNormal = lipgloss.NewStyle().
			Foreground(Text)
)
