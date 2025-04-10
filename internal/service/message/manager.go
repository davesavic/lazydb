package message

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davesavic/lazydb/internal/service/database"
)

// Manager is the manager used for triggering event messages.
type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

type NavigateDirectionMsg struct {
	Direction Direction
	Source    string
}

// NewNavigateDirectionCmd creates a new command for navigating in a direction.
// The directions are "up", "down", "left", and "right".
// The source is the name of the component that is sending the message.
func (m *Manager) NewNavigateDirectionCmd(direction Direction, source string) tea.Cmd {
	slog.Debug("NewNavigateDirectionCmd", "direction", direction, "source", source)
	return func() tea.Msg {
		return NavigateDirectionMsg{
			Direction: direction,
			Source:    source,
		}
	}
}

type ScreenName string

const (
	ScreenNameMain          ScreenName = "main"
	ScreenNameNewConnection ScreenName = "newConnection"
)

type ChangeScreenMsg struct {
	ScreenName ScreenName
}

// NewChangeScreenCmd creates a new command for changing the screen.
func (m *Manager) NewChangeScreenCmd(screenName ScreenName) tea.Cmd {
	slog.Debug("NewChangeScreenCmd", "screenName", screenName)
	return func() tea.Msg {
		return ChangeScreenMsg{
			ScreenName: screenName,
		}
	}
}

type PreviousScreenMsg struct{}

func (m *Manager) NewPreviousScreenCmd() tea.Cmd {
	slog.Debug("NewPreviousScreenCmd")
	return func() tea.Msg {
		return PreviousScreenMsg{}
	}
}

type NewAddConnectionMsg struct {
	Type     string
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (m *Manager) NewAddConnectionCmd(msg NewAddConnectionMsg) tea.Cmd {
	slog.Debug("NewAddConnectionCmd", "msg", msg)
	return func() tea.Msg {
		return msg
	}
}

type LoadConnectionMsg struct {
	Name string
}

func (m *Manager) NewLoadConnectionCmd(msg LoadConnectionMsg) tea.Cmd {
	slog.Debug("NewLoadConnectionCmd", "msg", msg)
	return func() tea.Msg {
		return msg
	}
}

type NewConnectionLoadedMsg struct{}

func (m *Manager) NewNewConnectionLoadedCmd() tea.Cmd {
	slog.Debug("NewNewConnectionLoadedCmd")
	return func() tea.Msg {
		return NewConnectionLoadedMsg{}
	}
}

type ExecuteQueryMsg struct {
	Query string
}

func (m *Manager) NewExecuteQueryCmd(query string) tea.Cmd {
	slog.Debug("NewExecuteQueryCmd", "query", query)
	return func() tea.Msg {
		return ExecuteQueryMsg{
			Query: query,
		}
	}
}

type QueryExecutedMsg struct {
	Result *database.QueryResult
}

func (m *Manager) NewQueryExecutedCmd(result *database.QueryResult) tea.Cmd {
	slog.Debug("NewQueryExecutedCmd", "result", result)
	return func() tea.Msg {
		return QueryExecutedMsg{
			Result: result,
		}
	}
}
