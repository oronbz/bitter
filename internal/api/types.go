package api

import "time"

// Build status codes from Bitrise API.
const (
	BuildStatusRunning = 0
	BuildStatusSuccess = 1
	BuildStatusFailed  = 2
	BuildStatusAborted = 3
)

type App struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	RepoURL   string `json:"repo_url"`
	RepoOwner string `json:"repo_owner"`
	RepoSlug  string `json:"repo_slug"`
	ProjectType string `json:"project_type"`
}

type Build struct {
	Slug                string     `json:"slug"`
	BuildNumber         int        `json:"build_number"`
	StatusText          string     `json:"status_text"`
	Status              int        `json:"status"`
	AbortReason         string     `json:"abort_reason"`
	Branch              string     `json:"branch"`
	CommitMessage       string     `json:"commit_message"`
	Tag                 string     `json:"tag"`
	TriggeredWorkflow   string     `json:"triggered_workflow"`
	TriggeredAt         time.Time  `json:"triggered_at"`
	StartedOnWorkerAt   *time.Time `json:"started_on_worker_at"`
	FinishedAt          *time.Time `json:"finished_at"`
	EnvironmentPrepareFinishedAt *time.Time `json:"environment_prepare_finished_at"`
	IsOnHold            bool       `json:"is_on_hold"`
	MachineTypeID       string     `json:"machine_type_id"`
}

type BuildLog struct {
	ExpiringRawLogURL string          `json:"expiring_raw_log_url"`
	IsArchived        bool            `json:"is_archived"`
	LogChunks         []LogChunk      `json:"log_chunks"`
	Timestamp         *time.Time      `json:"timestamp"`
}

type LogChunk struct {
	Chunk    string `json:"chunk"`
	Position int    `json:"position"`
}

type Pagination struct {
	TotalItemCount int    `json:"total_item_count"`
	PageItemLimit  int    `json:"page_item_limit"`
	Next           string `json:"next"`
}

type TriggerBuildParams struct {
	Branch     string `json:"branch"`
	WorkflowID string `json:"workflow_id"`
	CommitMessage string `json:"commit_message,omitempty"`
}

type AbortBuildParams struct {
	AbortReason           string `json:"abort_reason"`
	AbortWithSuccess      bool   `json:"abort_with_success"`
	SkipNotifications     bool   `json:"skip_notifications"`
}
