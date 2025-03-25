package newconnection

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/davesavic/lazydb/internal/message"
	"github.com/davesavic/lazydb/internal/ui/screen"
)

var _ screen.Screen = &NewConnection{}

type Connection struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type NewConnection struct {
	width  int
	height int

	form   *huh.Form
	result *Connection
}

func NewNewConnection() *NewConnection {
	result := &Connection{}

	return &NewConnection{
		form: huh.NewForm(
			huh.NewGroup(
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
func (n *NewConnection) Update(msg tea.Msg) (screen.Screen, tea.Cmd) {
	newForm, cmd := n.form.Update(msg)
	if f, ok := newForm.(*huh.Form); ok {
		n.form = f
	}

	if n.form.State == huh.StateCompleted {
		// Slog the form values
		slog.Debug("NewConnection.Update", "result", n.result)
		cmd = message.NewChangeScreenCmd(0)
	}

	return n, cmd
}

// View implements Screen.
func (n *NewConnection) View() string {
	return n.form.View()
}
