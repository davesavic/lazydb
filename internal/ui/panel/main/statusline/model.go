package statusline

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davesavic/lazydb/internal/message"
)

var _ tea.Model = &Model{}

type Model struct {
	status  string
	message string
	width   int
	height  int
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case message.StatusUpdateMsg:
		m.status = msg.Status
		m.message = msg.Message
	}

	return m, nil
}

// View implements tea.Model.
func (m *Model) View() string {
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 1).
		MarginRight(1).
		Height(1)

	statusMessage := lipgloss.NewStyle().Inherit(statusBarStyle)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusStyle.Render(m.status),
		statusMessage.Width(m.width-len(m.status)).Render(m.message),
	)

	return statusBarStyle.Width(m.width).MaxHeight(1).Render(bar)
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}
