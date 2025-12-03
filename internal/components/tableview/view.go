package tableview

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var placeholderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Italic(true)

func (m TableViewModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	if m.width == 0 || m.height == 0 {
		v.SetContent(placeholderStyle.Render("Table view (empty)"))
		return v
	}
	v.SetContent(m.table.View())
	return v
}
