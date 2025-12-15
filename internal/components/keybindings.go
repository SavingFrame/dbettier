package components

import (
	"charm.land/bubbles/v2/key"
)

// GlobalKeyMap defines keybindings available across all panes
type GlobalKeyMap struct {
	FocusLeft key.Binding
	FocusNext key.Binding
	FocusPrev key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// DefaultGlobalKeyMap returns the default global keybindings
var DefaultGlobalKeyMap = GlobalKeyMap{
	FocusLeft: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "focus tree"),
	),
	FocusNext: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "next pane"),
	),
	FocusPrev: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "prev pane"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.FocusLeft, k.FocusNext, k.FocusPrev},
		{k.Help, k.Quit},
	}
}
