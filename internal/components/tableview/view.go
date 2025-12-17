package tableview

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true)

	indicatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)
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
	return m.table.View() + "\n" + m.renderStatusBar()
}

func (m TableViewModel) renderStatusBar() string {
	var indicators []string

	// Table type indicator
	// if m.query.(type) == nil {
	if m.data.IsTableQuery() {
		indicators = append(indicators, " ")
	} else {
		indicators = append(indicators, " ")
	}

	// Vertical scroll indicator
	if m.table.GetHeight() > 2 {
		focusedRow, _ := m.table.FocusedPosition()
		totalRows := len(m.table.Rows())
		pageOffset := m.data.PageOffset()

		currentPos := focusedRow + 1 + pageOffset
		totalRowsStr := fmt.Sprintf("%d", totalRows+pageOffset)
		if m.data.CanFetchTotal() {
			totalRowsStr += "+"
		}
		indicators = append(indicators,
			indicatorStyle.Render(fmt.Sprintf("Row %d/%s", currentPos, totalRowsStr)))
	}

	// Horizontal scroll indicator
	if cols := m.table.Columns(); len(cols) > 1 {

		_, focusedCol := m.table.FocusedPosition()
		indicators = append(indicators,
			indicatorStyle.Render(fmt.Sprintf("Col %d/%d", focusedCol+1, len(cols))))
	}

	// Ordering indicator
	if orders := m.data.GetSortOrders(); len(orders) > 0 {
		var orderIndicators []string
		for _, orderCol := range orders {
			orderIndicators = append(orderIndicators, fmt.Sprintf("%s %s", orderCol.ColumnName, orderCol.Direction))
		}
		indicators = append(indicators,
			indicatorStyle.Render("Order: "+strings.Join(orderIndicators, ", ")))

	}

	if len(indicators) == 0 {
		return ""
	}

	statusBar := indicatorStyle.Render(strings.Join(indicators, " | "))

	// Add pagination message if present
	if msg := m.pagination.Message(); msg != "" {
		statusBar += "  " + messageStyle.Render(msg)
	}

	return statusBar
}

// formatNumber formats a number as a string.
func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}
