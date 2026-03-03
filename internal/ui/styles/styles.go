package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/api"
)

var (
	// Colors
	Amber   = lipgloss.Color("#FFBF00")
	DimGray = lipgloss.Color("#555555")
	White   = lipgloss.Color("#FFFFFF")
	Green   = lipgloss.Color("#00CC00")
	Red     = lipgloss.Color("#FF4444")
	Yellow  = lipgloss.Color("#FFDD00")
	Blue    = lipgloss.Color("#4488FF")

	// Panel borders
	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Amber)

	UnfocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DimGray)

	// Panel titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Amber).
			Padding(0, 1)

	// Build status indicators
	StatusSuccess = lipgloss.NewStyle().Foreground(Green).SetString("✓")
	StatusFailed  = lipgloss.NewStyle().Foreground(Red).SetString("✗")
	StatusRunning = lipgloss.NewStyle().Foreground(Yellow).SetString("●")
	StatusAborted = lipgloss.NewStyle().Foreground(DimGray).SetString("○")
	StatusOnHold  = lipgloss.NewStyle().Foreground(Blue).SetString("◆")

	// Build list columns
	BuildNumStyle      = lipgloss.NewStyle().Foreground(DimGray)
	BuildBranchStyle   = lipgloss.NewStyle().Foreground(Green)
	BuildWorkflowStyle = lipgloss.NewStyle().Foreground(Blue)
	BuildDurationStyle = lipgloss.NewStyle().Foreground(DimGray)
	BuildCommitStyle   = lipgloss.NewStyle().Foreground(White)

	// Error
	ErrorStyle = lipgloss.NewStyle().Foreground(Red)

	// Help overlay
	HelpOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Amber).
				Padding(1, 2)

	HelpKeyStyle    = lipgloss.NewStyle().Foreground(Amber).Bold(true).Width(16)
	HelpDescStyle   = lipgloss.NewStyle().Foreground(White)
	HelpHeaderStyle = lipgloss.NewStyle().Foreground(Amber).Bold(true).
			Underline(true).
			MarginBottom(1)

	// Dialog styles
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Amber).
			Padding(1, 2).
			Width(50)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Amber).
				MarginBottom(1)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().Foreground(Amber)
)

func BuildStatusIcon(status int, isOnHold bool) string {
	if isOnHold {
		return StatusOnHold.String()
	}
	switch status {
	case api.BuildStatusRunning:
		return StatusRunning.String()
	case api.BuildStatusSuccess:
		return StatusSuccess.String()
	case api.BuildStatusFailed:
		return StatusFailed.String()
	case api.BuildStatusAborted:
		return StatusAborted.String()
	default:
		return "?"
	}
}
