package buildlist

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Item struct {
	Build api.Build
}

func (i Item) Title() string       { return fmt.Sprintf("#%d %s", i.Build.BuildNumber, i.Build.Branch) }
func (i Item) Description() string { return i.Build.TriggeredWorkflow }
func (i Item) FilterValue() string {
	return fmt.Sprintf("%d %s %s %s", i.Build.BuildNumber, i.Build.Branch, i.Build.TriggeredWorkflow, i.Build.CommitMessage)
}

type Delegate struct {
	selectedStyle lipgloss.Style
	normalStyle   lipgloss.Style
}

func NewDelegate() Delegate {
	return Delegate{
		selectedStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(styles.Amber).
			PaddingLeft(1),
		normalStyle: lipgloss.NewStyle().
			PaddingLeft(2),
	}
}

func (d Delegate) Height() int                             { return 2 }
func (d Delegate) Spacing() int                            { return 0 }
func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d Delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	b := item.Build
	icon := styles.BuildStatusIcon(b.Status, b.IsOnHold)
	num := styles.BuildNumStyle.Render(fmt.Sprintf("#%d", b.BuildNumber))
	branch := styles.BuildBranchStyle.Render(b.Branch)
	commitMsg := firstLine(b.CommitMessage)
	if commitMsg == "" {
		commitMsg = "No commit message"
	}
	commit := styles.BuildCommitStyle.Render(commitMsg)
	workflow := styles.BuildWorkflowStyle.Render(b.TriggeredWorkflow)
	dur := styles.BuildDurationStyle.Render(formatDuration(b))

	line1 := fmt.Sprintf("%s %s", icon, commit)
	line2 := fmt.Sprintf("  %s %s %s %s", num, branch, workflow, dur)

	var style lipgloss.Style
	if index == m.Index() {
		style = d.selectedStyle
	} else {
		style = d.normalStyle
	}

	// Clamp output to list width so long branch names don't wrap and
	// break the panel height budget.
	style = style.MaxWidth(m.Width())

	content := style.Render(line1 + "\n" + line2)
	fmt.Fprint(w, content)
}

func firstLine(s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		return s[:i]
	}
	return s
}

func formatDuration(b api.Build) string {
	if b.StartedOnWorkerAt == nil {
		return "pending"
	}
	end := time.Now()
	if b.FinishedAt != nil {
		end = *b.FinishedAt
	}
	d := end.Sub(*b.StartedOnWorkerAt)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
}
