package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/messages"
)

func FetchLog(client *api.Client, appSlug, buildSlug string) tea.Cmd {
	return func() tea.Msg {
		log, err := client.GetBuildLog(appSlug, buildSlug)
		return messages.LogLoadedMsg{Log: log, Err: err}
	}
}
