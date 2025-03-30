package app

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/davesavic/lazydb/internal/service/database"
	"github.com/davesavic/lazydb/internal/service/message"
	screenmanager "github.com/davesavic/lazydb/internal/service/screen"
	"github.com/davesavic/lazydb/internal/ui/common"
)

var _ tea.Model = &App{}

// App is the main application struct that holds the state of the application.
type App struct {
	keys            *keybinding.Keymap
	screenManager   *screenmanager.Screen
	messageManager  *message.Manager
	configService   *config.Service
	databaseService *database.Postgres
}

func NewApp() *App {
	keys := keybinding.NewKeymap()
	db := database.NewPostgres()
	configService := config.NewService()

	return &App{
		keys:            keys,
		configService:   configService,
		databaseService: db,
		screenManager: screenmanager.NewScreen(&common.ScreenProps{
			MessageManager:  message.NewManager(),
			DatabaseService: db,
			ConfigService:   configService,
			Keymap:          keys,
		}),
	}
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	slog.Info("App.Init")
	return tea.Batch(
		a.screenManager.Init(),
	)
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

	case message.LoadConnectionMsg:
		slog.Debug("App.Update.LoadConnectionMsg", "msg", msg)
		consCfg, err := a.configService.GetConnection(msg.Name)
		if err != nil {
			// return a, a.messageManager.NewErrorCmd(err)
			slog.Error("App.Update.LoadConnectionMsg", "error", err)
			return a, tea.Quit
		}

		err = a.databaseService.Connect(*consCfg)
		if err != nil {
			// return a, a.messageManager.NewErrorCmd(err)
			slog.Error("App.Update.LoadConnectionMsg", "error", err)
			return a, tea.Quit
		}

		cmds = append(cmds, a.messageManager.NewNewConnectionLoadedCmd())
	case message.ExecuteQueryMsg:
		slog.Debug("App.Update.ExecuteQueryMsg", "msg", msg)
		result, err := a.databaseService.ExecuteQuery(msg.Query)
		if err != nil {
			// return a, a.messageManager.NewErrorCmd(err)
			slog.Error("App.Update.ExecuteQueryMsg", "error", err)
			return a, tea.Quit
		}

		cmds = append(cmds, a.messageManager.NewQueryExecutedCmd(result))
	}

	screenModel, cmd := a.screenManager.Update(msg)
	if sm, ok := screenModel.(*screenmanager.Screen); ok {
		a.screenManager = sm
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// View implements tea.Model.
func (a *App) View() string {
	return a.screenManager.View()
}
