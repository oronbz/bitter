package messages

import "github.com/oronbz/bitter/internal/api"

type AppsLoadedMsg struct {
	Apps []api.App
	Err  error
}

type BuildsLoadedMsg struct {
	Builds []api.Build
	Err    error
}

type LogLoadedMsg struct {
	Log string
	Err error
}

type BuildTriggeredMsg struct {
	Build api.Build
	Err   error
}

type BuildAbortedMsg struct {
	Err error
}

type BuildRefreshedMsg struct {
	Build api.Build
	Err   error
}

type TickMsg struct{}
