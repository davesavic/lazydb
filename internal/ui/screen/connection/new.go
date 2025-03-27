package connection

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/davesavic/lazydb/internal/service/message"
	"github.com/davesavic/lazydb/internal/ui/common"
)

type Connection struct {
	Type     string
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type NewConnection struct {
	width       int
	height      int
	screenProps *common.ScreenProps

	form   *huh.Form
	result *Connection
}

func NewNewConnection(props *common.ScreenProps) *NewConnection {
	result := &Connection{}

	return &NewConnection{
		screenProps: props,

		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Title("Type").Options(
					huh.Option[string]{Key: "postgres", Value: "Postgres"},
				).Value(&result.Type),
				huh.NewInput().Title("Host").Placeholder("localhost").Value(&result.Host),
				huh.NewInput().Title("Port").Placeholder("5432").Value(&result.Port),
				huh.NewInput().Title("Database").Placeholder("postgres").Value(&result.Database),
				huh.NewInput().Title("User").Placeholder("postgres").Value(&result.User),
				huh.NewInput().Title("Password").Placeholder("password").Value(&result.Password),
			),
		),

		result: result,
	}
}

// Init implements Screen.
func (n *NewConnection) Init() tea.Cmd {
	return n.form.Init()
}

// Update implements Screen.
func (n *NewConnection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, n.screenProps.Keymap.Cancel) {
			return n, n.screenProps.MessageManager.NewPreviousScreenCmd()
		}
	}

	newForm, cmd := n.form.Update(msg)
	if f, ok := newForm.(*huh.Form); ok {
		n.form = f
	}

	if n.form.State == huh.StateCompleted {
		slog.Debug("NewConnection.Update", "result", n.result)
		copiedResult := *n.result
		n.result = &Connection{}
		return n, n.screenProps.MessageManager.NewAddConnectionCmd(message.NewAddConnectionMsg{
			Type:     copiedResult.Type,
			Host:     copiedResult.Host,
			Port:     copiedResult.Port,
			Database: copiedResult.Database,
			User:     copiedResult.User,
			Password: copiedResult.Password,
		})
	}

	return n, cmd
}

// View implements Screen.
func (n *NewConnection) View() string {
	return n.form.View()
}
