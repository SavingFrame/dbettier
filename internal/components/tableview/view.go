package tableview

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	placeholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)
)

// View implements tea.Model interface
func (m TableViewModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}

// RenderContent returns the string representation of the view for composition
func (m TableViewModel) RenderContent() string {
	if !m.viewport.IsReady() {
		return placeholderStyle.Render("Table view (empty)")
	}
	return m.table.View() + "\n" + m.statusBar.View()
}
