package dbtree

import "charm.land/bubbles/v2/key"

// KeyMap defines keybindings for the database tree component
type KeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Left            key.Binding
	Right           key.Binding
	Space           key.Binding
	Search          key.Binding
	SearchNextMatch key.Binding
	SearchPrevMatch key.Binding
	ScrollUp        key.Binding
	ScrollDown      key.Binding
	Enter           key.Binding
	Quit            key.Binding
}

// DefaultKeyMap returns the default keybindings for the database tree
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/h", "collapse"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "expand"),
	),
	Space: key.NewBinding(
		key.WithKeys("space", " "),
		key.WithHelp("space", "toggle"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "page down"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "page up"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open table"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	SearchNextMatch: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next match"),
	),
	SearchPrevMatch: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "previous match"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Space, k.Enter},
		{k.ScrollUp, k.ScrollDown},
		{k.Quit},
	}
}
