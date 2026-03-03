package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/commands"
	"github.com/oronbz/bitter/internal/ui/components/applist"
	"github.com/oronbz/bitter/internal/ui/components/buildlist"
	"github.com/oronbz/bitter/internal/ui/components/dialog"
	"github.com/oronbz/bitter/internal/ui/components/helpoverlay"
	"github.com/oronbz/bitter/internal/ui/components/logview"
	"github.com/oronbz/bitter/internal/ui/components/statusbar"
	"github.com/oronbz/bitter/internal/ui/messages"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Panel int

const (
	PanelApps Panel = iota
	PanelBuilds
	PanelLogs
)

type Model struct {
	client *api.Client

	appList     applist.Model
	buildList   buildlist.Model
	logView     logview.Model
	statusBar   statusbar.Model
	helpOverlay helpoverlay.Model
	triggerDlg  dialog.TriggerModel
	confirmDlg  dialog.ConfirmModel
	spinner     spinner.Model

	focusedPanel  Panel
	selectedApp   *api.App
	selectedBuild *api.Build
	layout        Layout
	autoRefresh    bool
	pendingFocus   bool
	width         int
	height        int
	ready         bool
}

func NewModel(client *api.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	al := applist.New(30, 20)
	al.SetFocused(true)

	return Model{
		client:       client,
		appList:      al,
		buildList:    buildlist.New(40, 20),
		logView:      logview.New(50, 20),
		statusBar:    statusbar.New(),
		helpOverlay:  helpoverlay.New(),
		triggerDlg:   dialog.NewTrigger(),
		confirmDlg:   dialog.NewConfirm(),
		spinner:      s,
		focusedPanel: PanelApps,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		commands.FetchApps(m.client),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.resize()
		return m, nil

	case tea.MouseMsg:
		if m.triggerDlg.Visible() || m.confirmDlg.Visible() || m.helpOverlay.Visible() {
			break
		}
		target, ok := m.panelForMouse(msg.X, msg.Y)
		if !ok {
			break
		}
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			m.setFocus(target)
			return m, nil
		}
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			var cmd tea.Cmd
			switch target {
			case PanelApps:
				m.appList, cmd = m.appList.Update(msg)
			case PanelBuilds:
				m.buildList, cmd = m.buildList.Update(msg)
			case PanelLogs:
				m.logView, cmd = m.logView.Update(msg)
			}
			return m, cmd
		}

	case clearInfoMsg:
		m.statusBar.ClearInfo()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case messages.AppsLoadedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError(msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.appList.SetApps(msg.Apps)
		m.resize() // recalculate after items change pagination
		return m, nil

	case messages.BuildsLoadedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError(msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.buildList.SetBuilds(msg.Builds)
		m.resize() // recalculate after items change pagination
		if m.pendingFocus {
			m.pendingFocus = false
			m.cycleFocus(1)
		}
		m.autoRefresh = false
		for _, b := range msg.Builds {
			if b.Status == api.BuildStatusRunning && !b.IsOnHold {
				m.autoRefresh = true
				break
			}
		}
		if m.autoRefresh {
			cmds = append(cmds, commands.AutoRefreshTick())
		}
		return m, tea.Batch(cmds...)

	case messages.LogLoadedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError(msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.logView.SetContent(msg.Log)
		m.logView.ScrollToBottom()
		if m.pendingFocus {
			m.pendingFocus = false
			m.cycleFocus(1)
		}
		return m, nil

	case messages.BuildTriggeredMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Trigger failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		if m.selectedApp != nil {
			return m, commands.FetchBuilds(m.client, m.selectedApp.Slug)
		}
		return m, nil

	case messages.BuildAbortedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Abort failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		if m.selectedApp != nil {
			return m, commands.FetchBuilds(m.client, m.selectedApp.Slug)
		}
		return m, nil

	case messages.BuildRefreshedMsg:
		if msg.Err != nil {
			return m, nil
		}
		if m.selectedBuild != nil && m.selectedBuild.Slug == msg.Build.Slug {
			m.selectedBuild = &msg.Build
			if msg.Build.Status == api.BuildStatusRunning {
				cmds = append(cmds, commands.FetchLog(m.client, m.selectedApp.Slug, msg.Build.Slug))
			}
		}
		return m, tea.Batch(cmds...)

	case messages.TickMsg:
		if !m.autoRefresh || m.selectedApp == nil {
			return m, nil
		}
		batch := []tea.Cmd{commands.FetchBuilds(m.client, m.selectedApp.Slug)}
		if m.selectedBuild != nil && m.selectedBuild.Status == api.BuildStatusRunning {
			batch = append(batch,
				commands.FetchLog(m.client, m.selectedApp.Slug, m.selectedBuild.Slug),
			)
		}
		return m, tea.Batch(batch...)

	case dialog.TriggerSubmitMsg:
		if m.selectedApp == nil {
			return m, nil
		}
		m.statusBar.SetLoading("Triggering build...")
		return m, commands.TriggerBuild(m.client, m.selectedApp.Slug, api.TriggerBuildParams{
			Branch:     msg.Branch,
			WorkflowID: msg.Workflow,
		})

	case dialog.ConfirmYesMsg:
		if m.selectedApp == nil {
			return m, nil
		}
		switch msg.Action {
		case dialog.ConfirmAbort:
			if m.selectedBuild == nil {
				return m, nil
			}
			m.statusBar.SetLoading("Aborting build...")
			return m, commands.AbortBuild(m.client, m.selectedApp.Slug, m.selectedBuild.Slug)
		case dialog.ConfirmRebuild:
			if b, ok := m.buildList.SelectedBuild(); ok {
				m.statusBar.SetLoading("Triggering rebuild...")
				return m, commands.TriggerBuild(m.client, m.selectedApp.Slug, api.TriggerBuildParams{
					Branch:     b.Branch,
					WorkflowID: b.TriggeredWorkflow,
				})
			}
		}

	case dialog.ConfirmNoMsg:
		return m, nil
	}

	// Route to dialogs/overlays first if visible
	if m.triggerDlg.Visible() {
		var cmd tea.Cmd
		m.triggerDlg, cmd = m.triggerDlg.Update(msg)
		return m, cmd
	}
	if m.confirmDlg.Visible() {
		var cmd tea.Cmd
		m.confirmDlg, cmd = m.confirmDlg.Update(msg)
		return m, cmd
	}
	if m.helpOverlay.Visible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(keyMsg, Keys.Help) || key.Matches(keyMsg, Keys.Escape) || key.Matches(keyMsg, Keys.Quit) {
				m.helpOverlay.Toggle()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.helpOverlay, cmd = m.helpOverlay.Update(msg)
		return m, cmd
	}

	// Global keys (skip when a list is filtering, so keystrokes go to the filter input)
	filtering := m.appList.Filtering() || m.buildList.Filtering()
	if keyMsg, ok := msg.(tea.KeyMsg); ok && !filtering {
		switch {
		case key.Matches(keyMsg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(keyMsg, Keys.Help):
			m.helpOverlay.Toggle()
			return m, nil
		case key.Matches(keyMsg, Keys.Tab):
			m.cycleFocus(1)
			return m, nil
		case key.Matches(keyMsg, Keys.ShiftTab):
			m.cycleFocus(-1)
			return m, nil
		case key.Matches(keyMsg, Keys.Enter):
			return m, m.handleEnter()
		case key.Matches(keyMsg, Keys.Trigger):
			if m.selectedApp != nil {
				m.triggerDlg.Show()
			}
			return m, nil
		case key.Matches(keyMsg, Keys.Rebuild):
			if m.selectedApp != nil {
				if b, ok := m.buildList.SelectedBuild(); ok {
					m.confirmDlg.Show(
						fmt.Sprintf("Rebuild #%d (%s / %s)?", b.BuildNumber, b.Branch, b.TriggeredWorkflow),
						dialog.ConfirmRebuild,
					)
				}
			}
			return m, nil
		case key.Matches(keyMsg, Keys.Abort):
			if b, ok := m.buildList.SelectedBuild(); ok && b.Status == api.BuildStatusRunning {
				m.selectedBuild = &b
				m.confirmDlg.Show(fmt.Sprintf("Abort build #%d on %s?", b.BuildNumber, b.Branch), dialog.ConfirmAbort)
			}
			return m, nil
		case key.Matches(keyMsg, Keys.Refresh):
			return m, m.handleRefresh()
		case key.Matches(keyMsg, Keys.Open):
			openBrowser(m.buildURL())
			return m, nil
		case key.Matches(keyMsg, Keys.Copy):
			return m, m.copyURL()
		}
	}

	// Route to focused panel
	switch m.focusedPanel {
	case PanelApps:
		var cmd tea.Cmd
		m.appList, cmd = m.appList.Update(msg)
		cmds = append(cmds, cmd)
	case PanelBuilds:
		var cmd tea.Cmd
		m.buildList, cmd = m.buildList.Update(msg)
		cmds = append(cmds, cmd)
	case PanelLogs:
		var cmd tea.Cmd
		m.logView, cmd = m.logView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return m.spinner.View() + " Loading..."
	}

	// Early return for overlays — skip panel rendering
	if m.helpOverlay.Visible() {
		return m.helpOverlay.View()
	}
	if m.triggerDlg.Visible() {
		return m.triggerDlg.View()
	}
	if m.confirmDlg.Visible() {
		return m.confirmDlg.View()
	}

	l := m.layout
	appsPanel := m.renderPanel(m.appList.View(), l.AppsWidth, l.FullHeight, m.focusedPanel == PanelApps)

	var rightColumn string
	if m.selectedBuild != nil {
		buildsPanel := m.renderPanel(m.buildList.View(), l.RightWidth, l.BuildsHeight, m.focusedPanel == PanelBuilds)
		logPanel := m.renderPanel(m.logView.View(), l.RightWidth, l.LogHeight, m.focusedPanel == PanelLogs)
		rightColumn = lipgloss.JoinVertical(lipgloss.Left, buildsPanel, logPanel)
	} else {
		rightColumn = m.renderPanel(m.buildList.View(), l.RightWidth, l.FullHeight, m.focusedPanel == PanelBuilds)
	}

	panels := lipgloss.JoinHorizontal(lipgloss.Top, appsPanel, rightColumn)
	return lipgloss.JoinVertical(lipgloss.Left, panels, m.statusBar.View())
}

func (m Model) renderPanel(content string, width, height int, focused bool) string {
	style := styles.UnfocusedBorder
	if focused {
		style = styles.FocusedBorder
	}
	return style.Width(width).Height(height).MaxHeight(height + 2).Render(content)
}

func (m *Model) resize() {
	showLog := m.selectedBuild != nil
	m.layout = ComputeLayout(m.width, m.height, showLog)
	m.appList.SetSize(m.layout.AppsWidth, m.layout.FullHeight)
	if showLog {
		m.buildList.SetSize(m.layout.RightWidth, m.layout.BuildsHeight)
		m.logView.SetSize(m.layout.RightWidth, m.layout.LogHeight)
	} else {
		m.buildList.SetSize(m.layout.RightWidth, m.layout.FullHeight)
	}
	m.helpOverlay.SetSize(m.width, m.height)
	m.triggerDlg.SetSize(m.width, m.height)
	m.confirmDlg.SetSize(m.width, m.height)
}

func (m *Model) setFocus(panel Panel) {
	if m.focusedPanel == panel {
		return
	}
	m.appList.SetFocused(false)
	m.buildList.SetFocused(false)
	m.logView.SetFocused(false)

	m.focusedPanel = panel
	switch panel {
	case PanelApps:
		m.appList.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelApps)
	case PanelBuilds:
		m.buildList.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelBuilds)
	case PanelLogs:
		m.logView.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelLogs)
	}
}

func (m Model) panelForMouse(x, y int) (Panel, bool) {
	l := m.layout
	if y >= l.TotalHeight-statusBarHeight {
		return 0, false
	}
	if x < l.AppsWidth+2 {
		return PanelApps, true
	}
	if m.selectedBuild != nil && y >= l.BuildsHeight+2 {
		return PanelLogs, true
	}
	return PanelBuilds, true
}

func (m *Model) cycleFocus(dir int) {
	m.appList.SetFocused(false)
	m.buildList.SetFocused(false)
	m.logView.SetFocused(false)

	count := 2
	if m.selectedBuild != nil {
		count = 3
	}

	m.focusedPanel = Panel((int(m.focusedPanel) + dir + count) % count)

	switch m.focusedPanel {
	case PanelApps:
		m.appList.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelApps)
	case PanelBuilds:
		m.buildList.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelBuilds)
	case PanelLogs:
		m.logView.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelLogs)
	}
}

func (m *Model) handleEnter() tea.Cmd {
	switch m.focusedPanel {
	case PanelApps:
		if app, ok := m.appList.SelectedApp(); ok {
			m.selectedApp = &app
			m.selectedBuild = nil
			m.pendingFocus = true
			m.statusBar.SetLoading("Loading builds...")
			m.logView.SetContent("")
			m.resize()
			return commands.FetchBuilds(m.client, app.Slug)
		}
	case PanelBuilds:
		if build, ok := m.buildList.SelectedBuild(); ok {
			m.selectedBuild = &build
			m.pendingFocus = true
			m.resize()
			m.statusBar.SetLoading("Loading log...")
			return commands.FetchLog(m.client, m.selectedApp.Slug, build.Slug)
		}
	}
	return nil
}

func (m *Model) handleRefresh() tea.Cmd {
	switch m.focusedPanel {
	case PanelApps:
		m.statusBar.SetLoading("Refreshing apps...")
		return commands.FetchApps(m.client)
	case PanelBuilds:
		if m.selectedApp != nil {
			m.statusBar.SetLoading("Refreshing builds...")
			return commands.FetchBuilds(m.client, m.selectedApp.Slug)
		}
	case PanelLogs:
		if m.selectedApp != nil && m.selectedBuild != nil {
			m.statusBar.SetLoading("Refreshing log...")
			return commands.FetchLog(m.client, m.selectedApp.Slug, m.selectedBuild.Slug)
		}
	}
	return nil
}

func (m *Model) buildURL() string {
	if m.selectedApp == nil {
		return ""
	}
	// Prefer the currently highlighted build, fall back to the selected (entered) build.
	if b, ok := m.buildList.SelectedBuild(); ok {
		return fmt.Sprintf("https://app.bitrise.io/build/%s", b.Slug)
	}
	if m.selectedBuild != nil {
		return fmt.Sprintf("https://app.bitrise.io/build/%s", m.selectedBuild.Slug)
	}
	return ""
}

type clearInfoMsg struct{}

func (m *Model) copyURL() tea.Cmd {
	url := m.buildURL()
	if url == "" {
		return nil
	}
	if err := clipboard.WriteAll(url); err != nil {
		m.statusBar.SetError("Copy failed: " + err.Error())
		return nil
	}
	m.statusBar.SetInfo("Copied build URL to clipboard")
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} })
}

func openBrowser(url string) {
	if url == "" {
		return
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	if err := cmd.Start(); err == nil {
		go cmd.Wait()
	}
}
