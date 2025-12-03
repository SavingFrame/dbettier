package tableview

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var placeholderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Italic(true)

// RenderContent returns the string representation of the view for composition
func (m TableViewModel) RenderContent() string {
	if m.width == 0 || m.height == 0 {
		return placeholderStyle.Render("Table view (empty)")
	}
	return m.table.View()
}

// View implements tea.Model interface
func (m TableViewModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}
