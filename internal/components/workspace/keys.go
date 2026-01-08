package workspace

import "charm.land/bubbles/v2/key"

// KeyMap defines keybindings for the workspace tabs
type KeyMap struct {
	NextTab  key.Binding
	PrevTab  key.Binding
	CloseTab key.Binding
}

// DefaultKeyMap returns the default keybindings for the workspace
var DefaultKeyMap = KeyMap{
	NextTab: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "prev tab"),
	),
	CloseTab: key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "close tab"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.PrevTab, k.NextTab, k.CloseTab}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.PrevTab, k.NextTab, k.CloseTab},
	}
}
