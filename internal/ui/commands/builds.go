package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/messages"
)

func FetchBuilds(client *api.Client, appSlug string) tea.Cmd {
	return func() tea.Msg {
		builds, _, err := client.ListBuilds(appSlug, "")
		return messages.BuildsLoadedMsg{Builds: builds, Err: err}
	}
}

func TriggerBuild(client *api.Client, appSlug string, params api.TriggerBuildParams) tea.Cmd {
	return func() tea.Msg {
		build, err := client.TriggerBuild(appSlug, params)
		return messages.BuildTriggeredMsg{Build: build, Err: err}
	}
}

func AbortBuild(client *api.Client, appSlug, buildSlug string) tea.Cmd {
	return func() tea.Msg {
		err := client.AbortBuild(appSlug, buildSlug, api.AbortBuildParams{
			AbortReason: "Aborted via bitter TUI",
		})
		return messages.BuildAbortedMsg{Err: err}
	}
}
