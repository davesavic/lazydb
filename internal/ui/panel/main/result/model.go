package result

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/message"
)

var _ tea.Model = &Model{}

type Model struct {
	id     string
	keys   *keybinding.Keymap
	width  int
	height int
}

func NewModel(keys *keybinding.Keymap) *Model {
	return &Model{
		id:   "results",
		keys: keys,
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

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *Model) View() string {
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render("results panel")
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) Focus() {

}

func (m *Model) Blur() {

}
