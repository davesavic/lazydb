package result

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/service/message"
	"github.com/davesavic/lazydb/internal/service/plugin"
	"github.com/davesavic/lazydb/internal/ui/common"
	"github.com/evertras/bubble-table/table"
)

var _ tea.Model = &Model{}

type Model struct {
	id          string
	screenProps *common.ScreenProps
	width       int
	height      int

	table   *table.Model
	results *plugin.QueryResult
}

func NewModel(props *common.ScreenProps) *Model {
	return &Model{
		id:          "results",
		screenProps: props,
		table:       nil,
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
	case message.QueryExecutedMsg:
		m.results = msg.Result

		// Extract column names from the first row
		columns := make([]table.Column, 0)
		if len(m.results.Columns) > 0 {
			for _, title := range m.results.Columns {
				columns = append(columns, table.NewColumn(title, title, 25))
			}
		}

		rows := make([]table.Row, 0)
		for _, row := range m.results.Rows {
			rows = append(rows, table.NewRow(row))
		}

		t := table.New(columns).WithRows(rows)

		m.table = &t

	case tea.KeyMsg:
		switch {
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

	if m.table != nil {
		newTable, cmd := m.table.Update(msg)
		m.table = &newTable
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *Model) View() string {
	if m.table == nil {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Render("No results")
	}

	return lipgloss.NewStyle().
		Render(m.table.View())
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) Focus() {
	if m.table != nil {
		newTable := m.table.Focused(true)
		m.table = &newTable
	}
}

func (m *Model) Blur() {
	if m.table != nil {
		newTable := m.table.Focused(false)
		m.table = &newTable
	}
}
