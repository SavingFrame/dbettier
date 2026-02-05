package sqlcommandbarv2

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Execute key.Binding
}

var SQLCommandBarV2Keymap = KeyMap{
	Execute: key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "execute command"),
	),
}
