package tableview

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
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
	if !m.data.HasQuery() {
		return m.renderEmptyState()
	}
	return m.table.View() + "\n" + m.statusBar.View()
}

func (m TableViewModel) renderEmptyState() string {
	title := emptyStateTitleStyle().Render("󰆍 No query results yet")
	subtitle := emptyStateSubtitleStyle().Render("Run a query to fill this table.")

	sep := emptyStateHintStyle().Render(" ")
	hint1 := emptyStateBulletStyle().Render("•") + sep + emptyStateHintStyle().Render("Execute SQL with Alt+Enter in the editor")
	hint2 := emptyStateBulletStyle().Render("•") + sep + emptyStateHintStyle().Render("Press Enter on a table in the left tree")
	hint3 := emptyStateBulletStyle().Render("•") + sep + emptyStateHintStyle().Render("Press c on a table to open a query tab")

	content := strings.Join([]string{
		title,
		subtitle,
		"",
		hint1,
		hint2,
		hint3,
	}, "\n")

	contentWidth := lipgloss.Width(content)
	surfaceFill := lipgloss.NewStyle().Background(theme.Current().Colors.Surface)
	content = lipgloss.PlaceHorizontal(
		contentWidth,
		lipgloss.Left,
		content,
		lipgloss.WithWhitespaceStyle(surfaceFill),
	)

	card := emptyStateCardStyle().Render(content)
	return centerBlock(m.viewport.Width(), m.viewport.Height(), card)
}

func centerBlock(width, height int, content string) string {
	if width <= 0 || height <= 0 {
		return content
	}

	fill := lipgloss.NewStyle().Background(theme.Current().Colors.Base)
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceStyle(fill),
	)
}
