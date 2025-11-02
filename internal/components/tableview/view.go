package tableview

import (
	"github.com/charmbracelet/lipgloss"
)

var placeholderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Italic(true)

func (m TableViewModel) View() string {
	if m.width == 0 || m.height == 0 {
		return placeholderStyle.Render("Table view (empty)")
	}
	// Border is applied in root_screen.go, just return the table view
	return m.table.View()
}
