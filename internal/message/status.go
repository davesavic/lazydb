package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
