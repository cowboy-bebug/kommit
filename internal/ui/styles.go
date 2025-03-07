package ui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	MutedGray       = lipgloss.Color("#A1A1AA")
	MutedPurpleBlue = lipgloss.Color("#5A3FC0")
	NeuralGrey      = lipgloss.Color("#BDBDBD")
	SlateBlue       = lipgloss.Color("#64748B")
	SoftGreen       = lipgloss.Color("#6FCF97")
	WarmOrange      = lipgloss.Color("#F4A261")
	White           = lipgloss.Color("#FFFFFF")
)

// Styles
var (
	TitleStyle        = lipgloss.NewStyle()
	ItemStyle         = lipgloss.NewStyle().Padding(0, 1)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(White).Background(SlateBlue).Padding(0, 1)
	CheckedStyle      = lipgloss.NewStyle().Foreground(SoftGreen)
	UncheckedStyle    = lipgloss.NewStyle().Foreground(NeuralGrey)
	HelpStyle         = lipgloss.NewStyle().Foreground(MutedGray)
	KeyStyle          = lipgloss.NewStyle().Foreground(WarmOrange).Bold(true)
)
