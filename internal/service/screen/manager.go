package screenmanager

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davesavic/lazydb/internal/service/message"
	"github.com/davesavic/lazydb/internal/ui/common"
	"github.com/davesavic/lazydb/internal/ui/screen/connection"
	mainscreen "github.com/davesavic/lazydb/internal/ui/screen/main"
)

type ViewScreen interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

var _ tea.Model = &Screen{}

type Screen struct {
	messageManager *message.Manager
	screens        map[message.ScreenName]ViewScreen
	active         message.ScreenName
	previous       message.ScreenName
	width          int
	height         int
}

func NewScreen(props *common.ScreenProps) *Screen {
	screens := make(map[message.ScreenName]ViewScreen)

	screens[message.ScreenNameMain] = mainscreen.NewMain(props)
	screens[message.ScreenNameNewConnection] = connection.NewNewConnection(props)

	return &Screen{
		screens: screens,
		active:  message.ScreenNameMain,
	}
}

// Init implements tea.Model.
func (s *Screen) Init() tea.Cmd {
	return s.screens[s.active].Init()
}

// Update implements tea.Model.
func (s *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

		for i, sc := range s.screens {
			updatedScreen, cmd := sc.Update(msg)
			if updatedScreen, ok := updatedScreen.(ViewScreen); ok {
				s.screens[i] = updatedScreen
			}
			cmds = append(cmds, cmd)
		}

	case message.ChangeScreenMsg:
		s.previous = s.active
		s.active = msg.ScreenName

		cmds = append(cmds, s.screens[s.active].Init())

	case message.PreviousScreenMsg:
		if s.previous != "" {
			s.active = s.previous
		}
	}

	newScreen, cmd := s.screens[s.active].Update(msg)
	if newScreen, ok := newScreen.(ViewScreen); ok {
		s.screens[s.active] = newScreen
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

// View implements tea.Model.
func (s *Screen) View() string {
	return s.screens[s.active].View()
}
