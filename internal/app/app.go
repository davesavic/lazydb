package app

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/ui/screen"
)

var _ tea.Model = &App{}

// App is the main application struct that holds the state of the application.
type App struct {
	keys *keybinding.Keymap

	currentScreen screen.Screen
	screens       map[screen.Type]screen.Screen
}

func NewApp() *App {
	keys := keybinding.NewKeymap()
	screens := make(map[screen.Type]screen.Screen)
	screens[screen.TypeMain] = screen.NewMain(keys)

	return &App{
		keys:          keys,
		currentScreen: screens[screen.TypeMain],
		screens:       screens,
	}
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	slog.Info("App.Init")
	return nil
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keys.Help):
		}
	}

	newScreen, cmd := a.currentScreen.Update(msg)
	a.currentScreen = newScreen
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

// View implements tea.Model.
func (a *App) View() string {
	return a.currentScreen.View()
}
