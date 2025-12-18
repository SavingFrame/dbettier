package logpanel

import (
	"charm.land/bubbles/v2/key"
)

// KeyMap defines keybindings for the log panel
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Top      key.Binding
	Bottom   key.Binding
}

// DefaultKeyMap returns the default keybindings for log panel
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("k/up", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("j/down", "scroll down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn", "page down"),
	),
	Top: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "go to top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "go to bottom"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Top, k.Bottom},
	}
}
