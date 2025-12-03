package sqlcommandbar

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m SQLCommandBarModel) View() tea.View {
	var v tea.View
	v.AltScreen = true

	v.SetContent(fmt.Sprintf(
		"SQL Query:\n\n%s",
		m.textarea.View(),
	))
	return v
}
