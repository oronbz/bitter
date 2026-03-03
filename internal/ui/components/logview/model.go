package logview

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type Model struct {
	viewport viewport.Model
	focused  bool
}

func New(width, height int) Model {
	vp := viewport.New(width, height-1)
	vp.Style = lipgloss.NewStyle()
	return Model{viewport: vp}
}

func (m *Model) SetSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - 1
}

func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

func (m *Model) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) ScrollToBottom() {
	m.viewport.GotoBottom()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	title := styles.TitleStyle.Render("Log")
	return lipgloss.JoinVertical(lipgloss.Left, title, m.viewport.View())
}
