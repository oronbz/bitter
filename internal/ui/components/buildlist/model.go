package buildlist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Model struct {
	list    list.Model
	focused bool
}

func New(width, height int) Model {
	delegate := NewDelegate()
	l := list.New([]list.Item{}, delegate, width, height)
	l.Title = "Builds"
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

func (m *Model) SetBuilds(builds []api.Build) {
	items := make([]list.Item, len(builds))
	for i, b := range builds {
		items[i] = Item{Build: b}
	}
	m.list.SetItems(items)
}

func (m Model) SelectedBuild() (api.Build, bool) {
	item, ok := m.list.SelectedItem().(Item)
	if !ok {
		return api.Build{}, false
	}
	return item.Build, true
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
