package screen

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/ui"
	"github.com/davesavic/lazydb/internal/ui/connection"
	"github.com/davesavic/lazydb/internal/ui/query"
	"github.com/davesavic/lazydb/internal/ui/result"
	"github.com/davesavic/lazydb/internal/ui/statusline"
	"github.com/davesavic/lazydb/internal/ui/table"
)

var _ Screen = &Main{}

type Main struct {
	width  int
	height int

	activePanel PanelType
	navMap      NavigationMap

	// Panels
	connectionModel *connection.Model
	queryModel      *query.Model
	resultsModel    *result.Model
	tablesModel     *table.Model
	statusModel     *statusline.Model
}

func NewMain(keys *keybinding.Keymap) *Main {
	slog.Debug("NewMain")

	return &Main{
		activePanel:     PanelConnection,
		navMap:          NewNavigationMap(),
		connectionModel: connection.NewModel(keys),
		queryModel:      query.NewModel(keys),
		resultsModel:    result.NewModel(keys),
		tablesModel:     table.NewModel(keys),
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
func (m *Main) Update(msg tea.Msg) (Screen, tea.Cmd) {
	slog.Debug("Main.Update")
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeComponents(m.width, m.height)

	case ui.StatusUpdateMsg:
		newStatusPanel, cmd := m.statusModel.Update(msg)
		m.statusModel = newStatusPanel.(*statusline.Model)
		cmds = append(cmds, cmd)

	case ui.NavigationMsg:
		sourcePanel := ComponentIDToPane(msg.Source)

		if sourcePanel != m.activePanel {
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
	relationships map[PanelType]map[ui.NavigationDirection]PanelType
}

// NewNavigationMap creates a navigation map with default relationships
func NewNavigationMap() NavigationMap {
	nm := NavigationMap{
		relationships: make(map[PanelType]map[ui.NavigationDirection]PanelType),
	}

	// Initialize empty maps for each pane
	for i := range 4 {
		pane := PanelType(i)
		nm.relationships[pane] = make(map[ui.NavigationDirection]PanelType)
	}

	// Define the relationships

	// PanelConnection navigation
	nm.Set(PanelConnection, ui.NavDown, PanelTables)
	nm.Set(PanelConnection, ui.NavRight, PanelQuery)

	// PanelTables navigation
	nm.Set(PanelTables, ui.NavUp, PanelConnection)
	nm.Set(PanelTables, ui.NavRight, PanelResults)

	// PanelQuery navigation
	nm.Set(PanelQuery, ui.NavDown, PanelResults)
	nm.Set(PanelQuery, ui.NavLeft, PanelConnection)

	// PanelResults navigation
	nm.Set(PanelResults, ui.NavUp, PanelQuery)
	nm.Set(PanelResults, ui.NavLeft, PanelTables)

	return nm
}

// Set defines a directional relationship from one pane to another
func (nm *NavigationMap) Set(from PanelType, direction ui.NavigationDirection, to PanelType) {
	nm.relationships[from][direction] = to
}

// Navigate returns the target pane when navigating from a pane in a direction
func (nm *NavigationMap) Navigate(from PanelType, direction ui.NavigationDirection) (PanelType, bool) {
	if directions, exists := nm.relationships[from]; exists {
		if to, defined := directions[direction]; defined {
			return to, true
		}
	}
	return from, false // Return same pane if no navigation defined
}

// PaneToComponentID maps a pane type to its component ID
func PaneToComponentID(pane PanelType) string {
	switch pane {
	case PanelConnection:
		return "connections"
	case PanelTables:
		return "tables"
	case PanelQuery:
		return "query"
	case PanelResults:
		return "results"
	default:
		return ""
	}
}

// ComponentIDToPane maps a component ID to its pane type
func ComponentIDToPane(id string) PanelType {
	switch id {
	case "connections":
		return PanelConnection
	case "tables":
		return PanelTables
	case "query":
		return PanelQuery
	case "results":
		return PanelResults
	default:
		return PanelConnection
	}
}

type PanelType int

const (
	PanelConnection PanelType = iota
	PanelQuery
	PanelResults
	PanelTables
)
