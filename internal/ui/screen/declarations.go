package screen

import tea "github.com/charmbracelet/bubbletea"

type Type int

const (
	TypeMain Type = iota
	TypeNewConnection
)

type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
}
