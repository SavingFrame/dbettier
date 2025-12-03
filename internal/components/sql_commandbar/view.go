package sqlcommandbar

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// RenderContent returns the string representation of the view for composition
func (m SQLCommandBarModel) RenderContent() string {
	return fmt.Sprintf(
		"SQL Query:\n\n%s",
		m.textarea.View(),
	)
}

// View implements tea.Model interface
func (m SQLCommandBarModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}
