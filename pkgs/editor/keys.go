package editor

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Exit             key.Binding
	EnableInsertMode key.Binding
	InsertNextChar   key.Binding
	InsertNewLine    key.Binding
	Left             key.Binding
	Right            key.Binding
	Up               key.Binding
	Down             key.Binding
	Backspace        key.Binding
	Execute          key.Binding
	Space            key.Binding
	EndLineEdge      key.Binding
	StartLineEdge    key.Binding
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
	InsertNextChar: key.NewBinding(
		key.WithKeys("a"),
	),
	InsertNewLine: key.NewBinding(
		key.WithKeys("o"),
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
	StartLineEdge: key.NewBinding(
		key.WithKeys("^"),
	),
	EndLineEdge: key.NewBinding(
		key.WithKeys("$"),
	),
}

var InsertModeKeymap = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
	),
	Space: key.NewBinding(
		key.WithKeys("space"),
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
