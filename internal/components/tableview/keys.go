package tableview

import "charm.land/bubbles/v2/key"

// KeyMap defines keybindings for the table view component
type KeyMap struct {
	Enter        key.Binding
	Quit         key.Binding
	NextPage     key.Binding
	PreviousPage key.Binding
	Escape       key.Binding
}

// DefaultKeyMap returns the default keybindings for the table view
var DefaultKeyMap = KeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select row"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "bottom/next page"),
	),
	PreviousPage: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top/prev page"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextPage, k.PreviousPage, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPage, k.PreviousPage},
		{k.Quit},
	}
}
