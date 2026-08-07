package sqlcommandbar

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Execute key.Binding
}

var SQLCommandBarKeymap = KeyMap{
	Execute: key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "execute command"),
	),
}
