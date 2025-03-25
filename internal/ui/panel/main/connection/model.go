package connection

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/message"
)

var _ tea.Model = &Model{}

type Model struct {
	id     string
	keys   *keybinding.Keymap
	list   list.Model
	width  int
	height int
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
		case key.Matches(msg, m.keys.NavigateDown):
			cmds = append(cmds, message.RequestNavigationCmd(message.NavDown, m.id))
		case key.Matches(msg, m.keys.NavigateRight):
			cmds = append(cmds, message.RequestNavigationCmd(message.NavRight, m.id))
		case key.Matches(msg, m.keys.NavigateUp):
			cmds = append(cmds, message.RequestNavigationCmd(message.NavUp, m.id))
		case key.Matches(msg, m.keys.NavigateLeft):
			cmds = append(cmds, message.RequestNavigationCmd(message.NavLeft, m.id))
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

func NewModel(keys *keybinding.Keymap) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Connections"
	l.SetShowHelp(false)

	return &Model{
		id:   "connections",
		keys: keys,
		list: l,
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
