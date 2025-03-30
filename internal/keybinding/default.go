package keybinding

import "github.com/charmbracelet/bubbles/key"

// Keymap is a struct that holds the keybindings for the application.
type Keymap struct {
	// Global keybindings
	Quit          key.Binding
	Cancel        key.Binding
	Help          key.Binding
	NavigateUp    key.Binding
	NavigateDown  key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding

	// Query keybindings
	ExecuteQuery key.Binding

	// Connection keybindings
	AddConnection key.Binding
}

func NewKeymap() *Keymap {
	k := Keymap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
		NavigateUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "Navigate up"),
		),
		NavigateDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "Navigate down"),
		),
		NavigateLeft: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "Navigate left"),
		),
		NavigateRight: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "Navigate right"),
		),
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "Execute query"),
		),
		AddConnection: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "Add connection"),
		),
	}

	// apply user defined keybindings here
	// TODO: Load keymap bindings from a config file

	return &k
}

func (k Keymap) Bindings() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Help,
		k.NavigateUp,
		k.NavigateDown,
		k.NavigateLeft,
		k.NavigateRight,
		k.ExecuteQuery,
		k.AddConnection,
	}
}
