package applist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Item struct {
	App api.App
}

func (i Item) Title() string       { return i.App.Title }
func (i Item) Description() string { return i.App.RepoOwner + "/" + i.App.RepoSlug }
func (i Item) FilterValue() string { return i.App.Title }

type Model struct {
	list    list.Model
	focused bool
}

func New(width, height int) Model {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Amber).
		BorderLeftForeground(styles.Amber)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.DimGray).
		BorderLeftForeground(styles.Amber)

	l := list.New([]list.Item{}, delegate, width, height)
	l.Title = "Apps"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.Title = styles.TitleStyle
	l.DisableQuitKeybindings()

	return Model{list: l}
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

func (m *Model) SetApps(apps []api.App) {
	items := make([]list.Item, len(apps))
	for i, app := range apps {
		items[i] = Item{App: app}
	}
	m.list.SetItems(items)
}

func (m Model) SelectedApp() (api.App, bool) {
	item, ok := m.list.SelectedItem().(Item)
	if !ok {
		return api.App{}, false
	}
	return item.App, true
}

func (m Model) Filtering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
