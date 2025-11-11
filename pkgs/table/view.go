package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// View renders the table.
func (m Model) View() string {
	if len(m.cols) == 0 {
		return "No columns"
	}

	var s strings.Builder

	// Render header (with horizontal scrolling)
	s.WriteString(m.renderHeader())
	s.WriteString("\n")

	// Render rows (with vertical and horizontal scrolling)
	s.WriteString(m.renderRows())

	// Render scroll indicators if needed
	if m.focused {
		s.WriteString("\n")
		s.WriteString(m.renderScrollIndicators())
	}

	return s.String()
}

// renderHeader renders the table header (with horizontal scrolling).
func (m Model) renderHeader() string {
	var headers []string

	// Calculate visible columns based on width
	visibleCols := m.getVisibleColumns()

	for _, colIdx := range visibleCols {
		col := m.cols[colIdx]
		header := col.Title

		// Add sort indicator if this column is sorted
		sortIndicator := m.getSortIndicator(colIdx)
		if sortIndicator != "" {
			header = header + " " + sortIndicator
		}

		// Truncate or pad header to fit column width
		header = truncateOrPad(header, col.Width)

		// Apply header style
		style := m.styles.Header

		// Highlight header if this column is focused
		if m.focused && colIdx == m.focusedCol {
			style = style.Copy().Background(lipgloss.Color("62")).Bold(true)
		}

		headers = append(headers, style.Render(header))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}

// getSortIndicator returns the sort indicator for a column.
// Returns "↑" for ascending, "↓" for descending, with number prefix for multi-column sorts.
func (m Model) getSortIndicator(colIdx int) string {
	for i, sort := range m.orderColumns {
		if sort.ColumnIndex == colIdx {
			arrow := "↑"
			if sort.Direction == SortDesc {
				arrow = "↓"
			}

			// For multi-column sorts, show the order number
			if len(m.orderColumns) > 1 {
				// Use subscript numbers
				subscripts := []string{"₁", "₂", "₃", "₄", "₅", "₆", "₇", "₈", "₉"}
				if i < len(subscripts) {
					return subscripts[i] + arrow
				}
				return fmt.Sprintf("%d%s", i+1, arrow)
			}
			return arrow
		}
	}
	return ""
}

// renderRows renders the table rows (with scrolling).
func (m Model) renderRows() string {
	if len(m.rows) == 0 {
		return m.styles.Cell.Render("No data")
	}

	var rowStrings []string

	// Calculate visible rows based on height
	visibleRows := m.height - 2 // Subtract header and border
	if visibleRows <= 0 {
		visibleRows = len(m.rows)
	}

	// Calculate end index for rendering
	endIdx := m.scrollOffsetRow + visibleRows
	if endIdx > len(m.rows) {
		endIdx = len(m.rows)
	}

	// Render visible rows
	for rowIdx := m.scrollOffsetRow; rowIdx < endIdx; rowIdx++ {
		if rowIdx >= len(m.rows) {
			break
		}

		row := m.rows[rowIdx]
		rowStrings = append(rowStrings, m.renderRow(row, rowIdx))
	}

	return strings.Join(rowStrings, "\n")
}

// renderRow renders a single row (with horizontal scrolling).
func (m Model) renderRow(row Row, rowIdx int) string {
	var cells []string

	// Calculate visible columns
	visibleCols := m.getVisibleColumns()

	for _, colIdx := range visibleCols {
		col := m.cols[colIdx]
		var cellValue string

		// Get cell value (handle cases where row has fewer cells than columns)
		if colIdx < len(row) {
			cellValue = row[colIdx]
		} else {
			cellValue = ""
		}

		// Truncate or pad cell to fit column width
		cellValue = truncateOrPad(cellValue, col.Width)

		// Apply appropriate style based on focus
		style := m.getCellStyle(rowIdx, colIdx)

		cells = append(cells, style.Render(cellValue))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

// getCellStyle returns the appropriate style for a cell based on focus state.
func (m Model) getCellStyle(rowIdx, colIdx int) lipgloss.Style {
	if !m.focused {
		return m.styles.Cell
	}

	// Exact cell is focused - highest priority
	if rowIdx == m.focusedRow && colIdx == m.focusedCol {
		return m.styles.SelectedCell
	}

	// Cell is in the focused row
	if rowIdx == m.focusedRow {
		return m.styles.SelectedRow
	}

	// Cell is in the focused column
	if colIdx == m.focusedCol {
		return m.styles.SelectedCol
	}

	// Default cell style
	return m.styles.Cell
}

// truncateOrPad truncates or pads a string to the specified width.
// It handles multi-byte characters properly using runewidth.
func truncateOrPad(s string, width int) string {
	// Account for padding (1 space on each side)
	contentWidth := width - 2
	if contentWidth < 0 {
		contentWidth = 0
	}

	strWidth := runewidth.StringWidth(s)

	if strWidth > contentWidth {
		// Truncate with ellipsis
		if contentWidth <= 3 {
			return " " + runewidth.Truncate(s, contentWidth, "") + " "
		}
		return " " + runewidth.Truncate(s, contentWidth-3, "") + "..." + " "
	}

	// Pad with spaces
	padding := contentWidth - strWidth
	return " " + s + strings.Repeat(" ", padding) + " "
}

// getVisibleColumns returns the indices of columns that should be visible
// based on the current scroll offset and available width.
func (m Model) getVisibleColumns() []int {
	if m.width <= 0 || len(m.cols) == 0 {
		// If no width constraint, show all columns
		result := make([]int, len(m.cols))
		for i := range m.cols {
			result[i] = i
		}
		return result
	}

	var visibleCols []int
	currentWidth := 0

	// Start from scroll offset and add columns until we run out of width
	for i := m.scrollOffsetCol; i < len(m.cols); i++ {
		colWidth := m.cols[i].Width

		// Check if we have room for this column
		if currentWidth+colWidth > m.width && len(visibleCols) > 0 {
			break
		}

		visibleCols = append(visibleCols, i)
		currentWidth += colWidth
	}

	// If focused column is not visible, show at least the focused column
	if len(visibleCols) == 0 {
		visibleCols = []int{m.focusedCol}
	} else if visibleCols[0] > m.focusedCol || visibleCols[len(visibleCols)-1] < m.focusedCol {
		visibleCols = []int{m.focusedCol}
	}

	return visibleCols
}

// renderScrollIndicators renders scroll position indicators.
func (m Model) renderScrollIndicators() string {
	if len(m.rows) == 0 {
		return ""
	}

	var indicators []string

	// Vertical scroll indicator
	if m.height > 2 {
		visibleRows := m.height - 2
		totalRows := len(m.rows)

		if totalRows > visibleRows {
			currentPos := m.focusedRow + 1
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Row " + formatNumber(currentPos) + "/" + formatNumber(totalRows))
			indicators = append(indicators, indicator)
		}
	}

	// Horizontal scroll indicator
	if len(m.cols) > 0 {
		currentCol := m.focusedCol + 1
		totalCols := len(m.cols)

		if totalCols > 1 {
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Col " + formatNumber(currentCol) + "/" + formatNumber(totalCols))
			indicators = append(indicators, indicator)
		}
	}

	if len(indicators) == 0 {
		return ""
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Join(indicators, " | "))
}

// formatNumber formats a number as a string.
func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}
