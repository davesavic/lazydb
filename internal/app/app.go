package app

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/message"
	"github.com/davesavic/lazydb/internal/ui/common"
	"github.com/davesavic/lazydb/internal/ui/manager"
)

var _ tea.Model = &App{}

// App is the main application struct that holds the state of the application.
type App struct {
	keys           *keybinding.Keymap
	screenManager  *manager.Screen
	messageManager *message.Manager
}

func NewApp() *App {
	keys := keybinding.NewKeymap()

	return &App{
		keys: keys,
		screenManager: manager.NewScreen(&common.ScreenProps{
			MessageManager: message.NewManager(),
			Keymap:         keys,
		}),
	}
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	slog.Info("App.Init")
	return tea.Batch(a.screenManager.Init())
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		}
	}

	screenModel, cmd := a.screenManager.Update(msg)
	if sm, ok := screenModel.(*manager.Screen); ok {
		a.screenManager = sm
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// View implements tea.Model.
func (a *App) View() string {
	return a.screenManager.View()
}
