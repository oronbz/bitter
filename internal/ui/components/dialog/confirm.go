package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type ConfirmYesMsg struct{}
type ConfirmNoMsg struct{}

type ConfirmModel struct {
	message string
	visible bool
	width   int
	height  int
}

func NewConfirm() ConfirmModel {
	return ConfirmModel{}
}

func (m *ConfirmModel) Show(message string) {
	m.message = message
	m.visible = true
}

func (m *ConfirmModel) Hide() {
	m.visible = false
}

func (m ConfirmModel) Visible() bool {
	return m.visible
}

func (m *ConfirmModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.Hide()
			return m, func() tea.Msg { return ConfirmYesMsg{} }
		case "n", "N", "esc":
			m.Hide()
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	if !m.visible {
		return ""
	}

	title := styles.DialogTitleStyle.Render("Confirm")
	content := title + "\n\n" +
		m.message + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.DimGray).Render("y: yes  n/Esc: no")

	dialog := styles.DialogStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}
