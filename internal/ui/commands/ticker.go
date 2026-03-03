package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/ui/messages"
)

const RefreshInterval = 5 * time.Second

func AutoRefreshTick() tea.Cmd {
	return tea.Tick(RefreshInterval, func(time.Time) tea.Msg {
		return messages.TickMsg{}
	})
}
