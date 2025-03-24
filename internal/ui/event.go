package ui

import tea "github.com/charmbracelet/bubbletea"

type StatusUpdateMsg struct {
	Status  string
	Message string
}

func NewStatusUpdateCmd(status string, message string) tea.Cmd {
	return func() tea.Msg {
		return StatusUpdateMsg{
			Status:  status,
			Message: message,
		}
	}
}

type NavigationDirection int

const (
	NavUp NavigationDirection = iota
	NavDown
	NavLeft
	NavRight
)

type NavigationMsg struct {
	Direction NavigationDirection
	Source    string
}

func RequestNavigationCmd(direction NavigationDirection, source string) tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{
			Direction: direction,
			Source:    source,
		}
	}
}
