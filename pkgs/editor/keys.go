package editor

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Exit             key.Binding
	EnableInsertMode key.Binding
	Left             key.Binding
	Right            key.Binding
	Up               key.Binding
	Down             key.Binding
}

var NormalModeKeymap = KeyMap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q/ctrl+c", "exit"),
	),
	EnableInsertMode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert mode"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/left", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("r/right", "move right"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/down", "move down"),
	),
}

var InsertModeKeymap = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "move right"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Exit},
	}
}
