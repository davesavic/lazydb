package query

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/ui"
)

var _ tea.Model = &Model{}

type Model struct {
	id     string
	keys   *keybinding.Keymap
	width  int
	height int

	textarea textarea.Model
}

func NewModel(keys *keybinding.Keymap) *Model {
	textareaModel := textarea.New()
	textareaModel.Reset()
	textareaModel.SetCursor(0)

	return &Model{
		id:       "query",
		keys:     keys,
		textarea: textareaModel,
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
		case key.Matches(msg, m.keys.ExecuteQuery):
			// cmds = append(cmds, ui.RequestQueryCmd(m.id))
		case key.Matches(msg, m.keys.NavigateDown):
			cmds = append(cmds, ui.RequestNavigationCmd(ui.NavDown, m.id))
		case key.Matches(msg, m.keys.NavigateRight):
			cmds = append(cmds, ui.RequestNavigationCmd(ui.NavRight, m.id))
		case key.Matches(msg, m.keys.NavigateUp):
			cmds = append(cmds, ui.RequestNavigationCmd(ui.NavUp, m.id))
		case key.Matches(msg, m.keys.NavigateLeft):
			cmds = append(cmds, ui.RequestNavigationCmd(ui.NavLeft, m.id))
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
