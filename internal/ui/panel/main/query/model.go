package query

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/service/message"
	"github.com/davesavic/lazydb/internal/ui/common"
)

var _ tea.Model = &Model{}

type Model struct {
	id          string
	screenProps *common.ScreenProps
	width       int
	height      int

	textarea textarea.Model
}

func NewModel(props *common.ScreenProps) *Model {
	textareaModel := textarea.New()
	textareaModel.Reset()
	textareaModel.SetCursor(0)

	return &Model{
		id:          "query",
		screenProps: props,
		textarea:    textareaModel,
	}
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.screenProps.Keymap.ExecuteQuery):
			return m, m.screenProps.MessageManager.NewExecuteQueryCmd(m.textarea.Value())
		case key.Matches(msg, m.screenProps.Keymap.NavigateDown):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd(message.DirectionDown, m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateRight):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd(message.DirectionRight, m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateUp):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd(message.DirectionUp, m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateLeft):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd(message.DirectionLeft, m.id))
		}
	}

	newtextarea, cmd := m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.textarea = newtextarea

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *Model) View() string {
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(m.textarea.View())
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textarea.SetWidth(width)
	m.textarea.SetHeight(height)
}

func (m *Model) Focus() {
	m.textarea.Focus()
}

func (m *Model) Blur() {
	m.textarea.Blur()
}
