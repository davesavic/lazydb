package mainscreen

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/message"
	"github.com/davesavic/lazydb/internal/ui/common"
	"github.com/davesavic/lazydb/internal/ui/panel/main/connection"
	"github.com/davesavic/lazydb/internal/ui/panel/main/query"
	"github.com/davesavic/lazydb/internal/ui/panel/main/result"
	"github.com/davesavic/lazydb/internal/ui/panel/main/statusline"
	"github.com/davesavic/lazydb/internal/ui/panel/main/table"
)

type Main struct {
	width  int
	height int

	messageManager *message.Manager

	activePanel PanelID
	navMap      NavigationMap

	// Panels
	connectionModel *connection.Model
	queryModel      *query.Model
	resultsModel    *result.Model
	tablesModel     *table.Model
	statusModel     *statusline.Model
}

func NewMain(props *common.ScreenProps) *Main {
	slog.Debug("NewMain")

	return &Main{
		activePanel:     PanelConnection,
		navMap:          NewNavigationMap(),
		connectionModel: connection.NewModel(props),
		queryModel:      query.NewModel(props),
		resultsModel:    result.NewModel(props),
		tablesModel:     table.NewModel(props),
		statusModel:     statusline.NewModel(),
	}
}

// Init implements Screen.
func (m *Main) Init() tea.Cmd {
	slog.Debug("Main.Init")

	return tea.Batch(
		m.connectionModel.Init(),
		m.queryModel.Init(),
		m.resultsModel.Init(),
		m.tablesModel.Init(),
		m.statusModel.Init(),
	)
}

// Update implements Screen.
func (m *Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("Main.Update")
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeComponents(m.width, m.height)

	// case navigation.StatusUpdateMsg:
	// 	newStatusPanel, cmd := m.statusModel.Update(msg)
	// 	m.statusModel = newStatusPanel.(*statusline.Model)
	// 	cmds = append(cmds, cmd)

	case message.NavigateDirectionMsg:
		if PanelID(msg.Source) != m.activePanel {
			return m, nil
		}

		newPanel, ok := m.navMap.Navigate(m.activePanel, msg.Direction)
		if !ok {
			return m, nil
		}

		m.activePanel = newPanel

		m.connectionModel.Blur()
		m.queryModel.Blur()
		m.resultsModel.Blur()
		m.tablesModel.Blur()

		switch newPanel {
		case PanelConnection:
			m.connectionModel.Focus()
		case PanelQuery:
			m.queryModel.Focus()
		case PanelResults:
			m.resultsModel.Focus()
		case PanelTables:
			m.tablesModel.Focus()
		}
	}

	switch m.activePanel {
	case PanelConnection:
		newConnectionPanel, cmd := m.connectionModel.Update(msg)
		m.connectionModel = newConnectionPanel.(*connection.Model)
		cmds = append(cmds, cmd)

	case PanelQuery:
		newQueryPanel, cmd := m.queryModel.Update(msg)
		m.queryModel = newQueryPanel.(*query.Model)
		cmds = append(cmds, cmd)

	case PanelResults:
		newResultsPanel, cmd := m.resultsModel.Update(msg)
		m.resultsModel = newResultsPanel.(*result.Model)
		cmds = append(cmds, cmd)

	case PanelTables:
		newTablesPanel, cmd := m.tablesModel.Update(msg)
		m.tablesModel = newTablesPanel.(*table.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements Screen.
func (m *Main) View() string {
	return m.renderMainScreen()
}

func (m *Main) renderMainScreen() string {
	connectionSection := m.renderConnectionSection()
	querySection := m.renderQuerySection()

	mainView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		connectionSection,
		querySection,
	)

	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#4444444")).
		AlignHorizontal(lipgloss.Center).
		Width(m.width).UnsetMargins()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		statusStyle.Render(m.statusModel.View()),
	)
}

func (m *Main) stylePane(content string, active bool) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder())

	if active {
		style = style.BorderForeground(lipgloss.Color("#FF00FF"))
	}

	return style.Render(content)
}

func (m *Main) renderConnectionSection() string {
	connectionsView := m.stylePane(m.connectionModel.View(), m.activePanel == PanelConnection)
	tablesView := m.stylePane(m.tablesModel.View(), m.activePanel == PanelTables)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		connectionsView,
		tablesView,
	)
}

func (m *Main) renderQuerySection() string {
	queryView := m.stylePane(m.queryModel.View(), m.activePanel == PanelQuery)
	resultsView := m.stylePane(m.resultsModel.View(), m.activePanel == PanelResults)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		queryView,
		resultsView,
	)
}
func (m *Main) resizeComponents(width, height int) {
	// Add padding
	fullWidth := width

	width -= 4
	// height -= 5
	//
	// Reserve space for status bar
	contentHeight := height - 5

	// Left panel: 30% of width
	leftWidth := width * 3 / 10

	// Right panel: 70% of width
	rightWidth := width - leftWidth

	// Connections: 20% of content height
	connectionHeight := contentHeight * 2 / 10

	// Tables: Remaining left panel height
	tableHeight := contentHeight - connectionHeight

	// Query: 30% of right panel height
	queryHeight := contentHeight * 3 / 10

	// Results: Remaining right panel height
	resultsHeight := contentHeight - queryHeight

	// Update component dimensions
	m.connectionModel.SetSize(leftWidth, connectionHeight)
	m.tablesModel.SetSize(leftWidth, tableHeight)
	m.queryModel.SetSize(rightWidth, queryHeight)
	m.resultsModel.SetSize(rightWidth, resultsHeight)
	m.statusModel.SetSize(fullWidth, 1)
}

// Navigation map defines relationships between panes
type NavigationMap struct {
	// Maps each pane to its neighbors in each direction
	relationships map[PanelID]map[message.Direction]PanelID
}

// NewNavigationMap creates a navigation map with default relationships
func NewNavigationMap() NavigationMap {
	nm := NavigationMap{
		relationships: make(map[PanelID]map[message.Direction]PanelID),
	}

	nm.relationships[PanelConnection] = make(map[message.Direction]PanelID)
	nm.Set(PanelConnection, "down", PanelTables)
	nm.Set(PanelConnection, "right", PanelQuery)

	nm.relationships[PanelTables] = make(map[message.Direction]PanelID)
	nm.Set(PanelTables, "up", PanelConnection)
	nm.Set(PanelTables, "right", PanelResults)

	nm.relationships[PanelQuery] = make(map[message.Direction]PanelID)
	nm.Set(PanelQuery, "down", PanelResults)
	nm.Set(PanelQuery, "left", PanelConnection)

	nm.relationships[PanelResults] = make(map[message.Direction]PanelID)
	nm.Set(PanelResults, "up", PanelQuery)
	nm.Set(PanelResults, "left", PanelTables)

	return nm
}

// Set defines a directional relationship from one pane to another
func (nm *NavigationMap) Set(from PanelID, direction message.Direction, to PanelID) {
	nm.relationships[from][direction] = to
}

// Navigate returns the target pane when navigating from a pane in a direction
func (nm *NavigationMap) Navigate(from PanelID, direction message.Direction) (PanelID, bool) {
	if directions, exists := nm.relationships[from]; exists {
		if to, defined := directions[direction]; defined {
			return to, true
		}
	}
	return from, false
}

type PanelID string

const (
	PanelConnection PanelID = "connections"
	PanelQuery      PanelID = "query"
	PanelResults    PanelID = "results"
	PanelTables     PanelID = "tables"
)
