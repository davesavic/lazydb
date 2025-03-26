package connection

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/message"
	"github.com/davesavic/lazydb/internal/ui/common"
)

var _ tea.Model = &Model{}

type Model struct {
	id          string
	screenProps *common.ScreenProps
	list        list.Model
	width       int
	height      int
}

type listItem struct {
	name        string
	description string
}

func (l listItem) Title() string {
	return l.name
}

func (l listItem) Description() string {
	return l.description
}

func (l listItem) FilterValue() string {
	return l.name
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
		case msg.String() == "q":
			return m, nil
		case key.Matches(msg, m.screenProps.Keymap.AddConnection):
			cmds = append(cmds, m.screenProps.MessageManager.NewChangeScreenCmd(message.ScreenNameNewConnection))
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

	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *Model) View() string {
	return lipgloss.
		NewStyle().
		Width(m.width).
		Height(m.height).
		Render(m.list.View())
}

func NewModel(props *common.ScreenProps) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Connections"
	l.SetShowHelp(false)

	return &Model{
		id:          "connections",
		screenProps: props,
		list:        l,
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func (m *Model) Focus() {

}

func (m *Model) Blur() {

}
