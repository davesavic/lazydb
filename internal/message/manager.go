package message

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

// Manager is the manager used for triggering event messages.
type Manager struct {
}

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
