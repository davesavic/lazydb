package message

import tea "github.com/charmbracelet/bubbletea"

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
