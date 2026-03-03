package dialog

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/bitter/internal/ui/styles"
)

type TriggerSubmitMsg struct {
	Branch   string
	Workflow string
}

type TriggerModel struct {
	branchInput   textinput.Model
	workflowInput textinput.Model
	focusIndex    int
	visible       bool
	width         int
	height        int
}

func NewTrigger() TriggerModel {
	bi := textinput.New()
	bi.Placeholder = "main"
	bi.CharLimit = 128
	bi.Width = 40
	bi.Prompt = "Branch: "
	bi.PromptStyle = lipgloss.NewStyle().Foreground(styles.Amber)

	wi := textinput.New()
	wi.Placeholder = "primary"
	wi.CharLimit = 128
	wi.Width = 40
	wi.Prompt = "Workflow: "
	wi.PromptStyle = lipgloss.NewStyle().Foreground(styles.Amber)

	return TriggerModel{
		branchInput:   bi,
		workflowInput: wi,
	}
}

func (m *TriggerModel) Show() {
	m.visible = true
	m.focusIndex = 0
	m.branchInput.SetValue("")
	m.workflowInput.SetValue("")
	m.branchInput.Focus()
	m.workflowInput.Blur()
}

func (m *TriggerModel) Hide() {
	m.visible = false
	m.branchInput.Blur()
	m.workflowInput.Blur()
}

func (m TriggerModel) Visible() bool {
	return m.visible
}

func (m *TriggerModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m TriggerModel) Update(msg tea.Msg) (TriggerModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Hide()
			return m, nil
		case "tab", "shift+tab":
			m.focusIndex = (m.focusIndex + 1) % 2
			if m.focusIndex == 0 {
				m.branchInput.Focus()
				m.workflowInput.Blur()
			} else {
				m.branchInput.Blur()
				m.workflowInput.Focus()
			}
			return m, nil
		case "enter":
			branch := m.branchInput.Value()
			if branch == "" {
				branch = "main"
			}
			workflow := m.workflowInput.Value()
			if workflow == "" {
				workflow = "primary"
			}
			m.Hide()
			return m, func() tea.Msg {
				return TriggerSubmitMsg{Branch: branch, Workflow: workflow}
			}
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.branchInput, cmd = m.branchInput.Update(msg)
	} else {
		m.workflowInput, cmd = m.workflowInput.Update(msg)
	}
	return m, cmd
}

func (m TriggerModel) View() string {
	if !m.visible {
		return ""
	}

	title := styles.DialogTitleStyle.Render("Trigger Build")
	content := title + "\n\n" +
		m.branchInput.View() + "\n\n" +
		m.workflowInput.View() + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.DimGray).Render("Enter: submit  Tab: next field  Esc: cancel")

	dialog := styles.DialogStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}
