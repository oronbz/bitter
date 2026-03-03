package statusbar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Panel int

const (
	PanelApps Panel = iota
	PanelBuilds
	PanelLogs
)

type Model struct {
	panel   Panel
	errMsg  string
	loading string
}

func New() Model {
	return Model{}
}

func (m *Model) SetPanel(panel Panel) {
	m.panel = panel
}

func (m *Model) SetError(msg string) {
	m.errMsg = msg
}

func (m *Model) ClearError() {
	m.errMsg = ""
}

func (m *Model) SetLoading(msg string) {
	m.loading = msg
}

func (m *Model) ClearLoading() {
	m.loading = ""
}

type hint struct {
	key  string
	desc string
}

var (
	brandStyle = lipgloss.NewStyle().
			Background(styles.Amber).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)
	keyStyle = lipgloss.NewStyle().Foreground(styles.Amber).Bold(true)
	dimStyle = lipgloss.NewStyle().Foreground(styles.DimGray)
	errStyle = lipgloss.NewStyle().Foreground(styles.Red)

	appsHints = []hint{
		{"↑/k", "up"}, {"↓/j", "down"}, {"⏎", "select"},
		{"Tab", "panel"}, {"/", "filter"}, {"?", "help"}, {"q", "quit"},
	}
	buildsHints = []hint{
		{"↑/k", "up"}, {"↓/j", "down"}, {"⏎", "select"},
		{"Tab", "panel"}, {"t", "trigger"}, {"a", "abort"},
		{"/", "filter"}, {"?", "help"}, {"q", "quit"},
	}
	logsHints = []hint{
		{"↑/k", "up"}, {"↓/j", "down"}, {"C-d", "pgdn"}, {"C-u", "pgup"},
		{"Tab", "panel"}, {"r", "refresh"}, {"?", "help"}, {"q", "quit"},
	}
)

func (m Model) View() string {
	brand := brandStyle.Render("BITTER")

	var hints []hint
	switch m.panel {
	case PanelApps:
		hints = appsHints
	case PanelBuilds:
		hints = buildsHints
	case PanelLogs:
		hints = logsHints
	}

	var parts []string
	for _, h := range hints {
		parts = append(parts, keyStyle.Render(h.key)+" "+h.desc)
	}
	hintsStr := strings.Join(parts, "  ")

	var rightText string
	if m.errMsg != "" {
		rightText = "  " + errStyle.Render(m.errMsg)
	} else if m.loading != "" {
		rightText = "  " + dimStyle.Render(m.loading)
	}

	return brand + " " + dimStyle.Render("│") + " " + hintsStr + rightText
}
