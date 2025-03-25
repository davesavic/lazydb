package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ChangeScreenMsg struct {
	ScreenType int
}

func NewChangeScreenCmd(screenType int) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenMsg{
			ScreenType: screenType,
		}
	}
}
