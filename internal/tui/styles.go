package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("86")
	ColorSecondary = lipgloss.Color("241")
	ColorSuccess   = lipgloss.Color("82")
	ColorError     = lipgloss.Color("196")
	ColorWarning   = lipgloss.Color("214")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginTop(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2)

	CommitSHAStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	CommitDateStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	CommitAuthorStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary)

	RepoHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)
)
