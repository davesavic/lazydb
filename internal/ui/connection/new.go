package connection

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var _ tea.Model = &NewConnectionModel{}

type NewConnectionModel struct {
	name string

	form *huh.Form
}

func NewNewConnectionModel() *NewConnectionModel {
	n := &NewConnectionModel{}
	n.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").CharLimit(50).Value(&n.name),
		),
	)

	return n
}

// Init implements tea.Model.
func (n *NewConnectionModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (n *NewConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	newForm, cmd := n.form.Update(msg)
	if f, ok := newForm.(*huh.Form); ok {
		n.form = f
		cmds = append(cmds, cmd)
	}

	return n, tea.Batch(cmds...)
}

// View implements tea.Model.
func (n *NewConnectionModel) View() string {
	return n.form.View()
}
