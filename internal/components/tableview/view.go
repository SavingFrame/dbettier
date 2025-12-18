package tableview

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// View implements tea.Model interface
func (m TableViewModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.SetContent(m.RenderContent())
	return v
}

// RenderContent returns the string representation of the view for composition
func (m TableViewModel) RenderContent() string {
	if m.isLoading {
		return fmt.Sprintf("\n\n   %s Fetching data...\n\n", m.spinner.View())
	}
	if !m.viewport.IsReady() {
		return placeholderStyle().Render("Table view (empty)")
	}
	return m.table.View() + "\n" + m.statusBar.View()
}
