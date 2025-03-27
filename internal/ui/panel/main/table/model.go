package table

import (
	"log/slog"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/service/message"
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
func NewModel(props *common.ScreenProps) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Tables"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return &Model{
		id:          "tables",
		screenProps: props,
		list:        l,
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
	case message.NewConnectionLoadedMsg:
		tables, err := m.screenProps.DatabaseService.GetTables()
		if err != nil {
			slog.Error("failed to get tables", "error", err)
			// return m, m.screenProps.MessageManager.NewErrorCmd(err)
		}

		sort.Strings(tables)

		items := make([]list.Item, 0, len(tables))
		for _, t := range tables {
			items = append(items, listItem{name: t})
		}

		m.list.SetItems(items)

	case tea.KeyMsg:
		switch {
		case msg.String() == "q":
			return m, nil
		case key.Matches(msg, m.screenProps.Keymap.NavigateDown):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd("down", m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateRight):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd("right", m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateUp):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd("up", m.id))
		case key.Matches(msg, m.screenProps.Keymap.NavigateLeft):
			cmds = append(cmds, m.screenProps.MessageManager.NewNavigateDirectionCmd("left", m.id))
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

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func (m *Model) Focus() {

}

func (m *Model) Blur() {

}
