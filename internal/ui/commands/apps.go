package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/messages"
)

func FetchApps(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		apps, _, err := client.ListApps("")
		return messages.AppsLoadedMsg{Apps: apps, Err: err}
	}
}
