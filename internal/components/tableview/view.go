package tableview

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/pkgs/table"
)

var placeholderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Italic(true)

// RenderContent returns the string representation of the view for composition
func (m TableViewModel) RenderContent() string {
	if m.width == 0 || m.height == 0 {
		return placeholderStyle.Render("Table view (empty)")
	}
	return m.table.View() + "\n" + renderScrollIndicators(m.table, m)
}

// View implements tea.Model interface
func (m TableViewModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}

func renderScrollIndicators(t table.Model, m TableViewModel) string {
	if len(t.Rows()) == 0 {
		return ""
	}

	var indicators []string

	// Table type indicator
	// if m.query.(type) == nil {
	if _, ok := m.query.(*sharedcomponents.TableQuery); ok {
		indicators = append(indicators, " ")
	} else {
		indicators = append(indicators, " ")
	}

	focusedRow, focusedCol := t.FocusedPosition()
	// Vertical scroll indicator
	if t.GetHeight() > 2 {
		totalRows := len(t.Rows())

		var pageOffset int
		if m.query != nil {
			pageOffset = m.query.PageOffset()
		}
		currentPos := focusedRow + 1 + pageOffset
		totalRowsString := formatNumber(totalRows + pageOffset)
		if m.canFetchTotal {
			totalRowsString += "+"
		}
		indicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Row " + formatNumber(currentPos) + "/" + totalRowsString)
		indicators = append(indicators, indicator)
	}

	// Horizontal scroll indicator
	if len(t.Columns()) > 0 {
		currentCol := focusedCol + 1
		totalCols := len(t.Columns())

		if totalCols > 1 {
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Col " + formatNumber(currentCol) + "/" + formatNumber(totalCols))
			indicators = append(indicators, indicator)
		}
	}

	// Ordering indicator
	if len(t.OrderColumns()) > 0 {
		var orderIndicators []string
		for _, orderCol := range m.query.GetSortOrders() {
			orderIndicators = append(orderIndicators, fmt.Sprintf("%s %s", orderCol.ColumnName, orderCol.Direction))
		}
		indicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Order: " + strings.Join(orderIndicators, ", "))
		indicators = append(indicators, indicator)

	}

	if len(indicators) == 0 {
		return ""
	}

	leftSection := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Join(indicators, " | "))

	// If there's a custom message, display it in the center
	if m.customMessage != "" {
		centerSection := lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true).
			Render(m.customMessage)

		return leftSection + "  " + centerSection
	}

	return leftSection
}

// formatNumber formats a number as a string.
func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}
