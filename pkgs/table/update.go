package table

import (
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// SearchExitMsg is sent when search mode is exited.
type SearchExitMsg struct{}

// SearchUpdateMsg is sent when search query changes.
type SearchUpdateMsg struct {
	Query   string
	Matches int
}

// Update handles messages and updates the table state.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			return m.handleSearchInput(msg)
		}

		switch msg.String() {
		// Enter search mode
		case "/":
			m.searchMode = true
			m.searchQuery = ""
			m.searchMatches = nil
			m.searchMatchIndex = -1
			return m, nil

		// Navigate to next search match
		case "n":
			if len(m.searchMatches) > 0 {
				m.nextSearchMatch()
			}
			return m, nil

		// Navigate to previous search match
		case "N":
			if len(m.searchMatches) > 0 {
				m.prevSearchMatch()
			}
			return m, nil

		// Row navigation (up/down) - vim keys
		case "k":
			m.moveUp()
		case "j":
			m.moveDown()

		// Column navigation (left/right) - vim keys
		case "h":
			m.moveLeft()
		case "l":
			m.moveRight()

		// Jump to start/end - vim keys
		case "g":
			m.ScrollToTop()
		case "G":
			m.ScrollToBottom()

		// Page navigation - vim keys
		case "ctrl+b":
			m.pageUp()
		case "ctrl+f":
			m.pageDown()
		case "ctrl+u":
			m.halfPageUp()
		case "ctrl+d":
			m.halfPageDown()

		// Yank (copy) cell value - vim key
		case "y":
			return m, m.yankCell()

		// Sort by focused column - 's' key
		case "s":
			return m.toggleSort()

		// Clear all sorts - Shift+S
		case "S":
			return m.clearSort()
		}
	}

	return m, nil
}

// handleSearchInput handles key input during search mode.
func (m Model) handleSearchInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Exit search mode and clear highlights
		m.searchMode = false
		m.searchQuery = ""
		m.searchMatches = nil
		m.searchMatchIndex = -1
		return m, func() tea.Msg { return SearchExitMsg{} }

	case "enter":
		// Confirm search and exit search mode (keep highlights until next search)
		m.searchMode = false
		if len(m.searchMatches) > 0 && m.searchMatchIndex >= 0 {
			// Jump to current match
			match := m.searchMatches[m.searchMatchIndex]
			m.focusedRow = match.Row
			m.focusedCol = match.Col
			m.updateScrollRow()
			m.updateScrollCol()
		}
		return m, nil

	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.updateSearchMatches()
		}
		return m, func() tea.Msg {
			return SearchUpdateMsg{Query: m.searchQuery, Matches: len(m.searchMatches)}
		}

	default:
		// Add character to search query (only printable characters)
		if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] < 127 {
			m.searchQuery += msg.String()
			m.updateSearchMatches()
			return m, func() tea.Msg {
				return SearchUpdateMsg{Query: m.searchQuery, Matches: len(m.searchMatches)}
			}
		}
	}

	return m, nil
}

// updateSearchMatches finds all cells matching the search query.
func (m *Model) updateSearchMatches() {
	m.searchMatches = nil
	m.searchMatchIndex = -1

	if m.searchQuery == "" {
		return
	}

	query := strings.ToLower(m.searchQuery)

	for rowIdx, row := range m.rows {
		for colIdx, cell := range row {
			if strings.Contains(strings.ToLower(cell), query) {
				m.searchMatches = append(m.searchMatches, SearchMatch{
					Row: rowIdx,
					Col: colIdx,
				})
			}
		}
	}

	// If we have matches, set index to first match
	if len(m.searchMatches) > 0 {
		m.searchMatchIndex = 0
		// Jump to first match
		match := m.searchMatches[0]
		m.focusedRow = match.Row
		m.focusedCol = match.Col
		m.updateScrollRow()
		m.updateScrollCol()
	}
}

// nextSearchMatch moves to the next search match.
func (m *Model) nextSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}

	m.searchMatchIndex = (m.searchMatchIndex + 1) % len(m.searchMatches)
	match := m.searchMatches[m.searchMatchIndex]
	m.focusedRow = match.Row
	m.focusedCol = match.Col
	m.updateScrollRow()
	m.updateScrollCol()
}

// prevSearchMatch moves to the previous search match.
func (m *Model) prevSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}

	m.searchMatchIndex--
	if m.searchMatchIndex < 0 {
		m.searchMatchIndex = len(m.searchMatches) - 1
	}
	match := m.searchMatches[m.searchMatchIndex]
	m.focusedRow = match.Row
	m.focusedCol = match.Col
	m.updateScrollRow()
	m.updateScrollCol()
}

// ClearSearch clears the search state.
func (m *Model) ClearSearch() {
	m.searchMode = false
	m.searchQuery = ""
	m.searchMatches = nil
	m.searchMatchIndex = -1
}

// moveUp moves the focus up one row.
func (m *Model) moveUp() {
	if m.focusedRow > 0 {
		m.focusedRow--
		m.updateScrollRow()
	}
}

// moveDown moves the focus down one row.
func (m *Model) moveDown() {
	if m.focusedRow < len(m.rows)-1 {
		m.focusedRow++
		m.updateScrollRow()
	}
}

// moveLeft moves the focus left one column.
func (m *Model) moveLeft() {
	if m.focusedCol > 0 {
		m.focusedCol--
		m.updateScrollCol()
	}
}

// moveRight moves the focus right one column.
func (m *Model) moveRight() {
	if m.focusedCol < len(m.cols)-1 {
		m.focusedCol++
		m.updateScrollCol()
	}
}

// pageUp moves the focus up one page (based on visible height).
func (m *Model) pageUp() {
	if m.height <= 2 {
		m.moveUp()
		return
	}

	visibleRows := m.height - 2 // Subtract header and border
	m.focusedRow -= visibleRows
	if m.focusedRow < 0 {
		m.focusedRow = 0
	}
	m.updateScrollRow()
}

// pageDown moves the focus down one page (based on visible height).
func (m *Model) pageDown() {
	if m.height <= 2 {
		m.moveDown()
		return
	}

	visibleRows := m.height - 2 // Subtract header and border
	m.focusedRow += visibleRows
	if m.focusedRow >= len(m.rows) {
		m.focusedRow = len(m.rows) - 1
	}
	m.updateScrollRow()
}

// updateScrollRow adjusts the scroll offset to keep the focused row visible.
func (m *Model) updateScrollRow() {
	if m.height <= 2 {
		return
	}

	visibleRows := m.height - 2 // Subtract header and border

	// If focused row is above the visible area, scroll up
	if m.focusedRow < m.scrollOffsetRow {
		m.scrollOffsetRow = m.focusedRow
	}

	// If focused row is below the visible area, scroll down
	if m.focusedRow >= m.scrollOffsetRow+visibleRows {
		m.scrollOffsetRow = m.focusedRow - visibleRows + 1
	}

	// Ensure scroll offset doesn't go negative
	if m.scrollOffsetRow < 0 {
		m.scrollOffsetRow = 0
	}
}

// updateScrollCol adjusts the scroll offset to keep the focused column visible.
func (m *Model) updateScrollCol() {
	if m.width <= 0 || len(m.cols) == 0 {
		return
	}

	// If focused column is before scroll offset, scroll left
	if m.focusedCol < m.scrollOffsetCol {
		m.scrollOffsetCol = m.focusedCol
		return
	}

	// Calculate how many columns fit in the available width
	currentWidth := 0
	visibleColCount := 0

	for i := m.scrollOffsetCol; i < len(m.cols); i++ {
		colWidth := m.cols[i].Width
		if currentWidth+colWidth > m.width && visibleColCount > 0 {
			break
		}
		currentWidth += colWidth
		visibleColCount++

		// If we've reached the focused column, we're done
		if i == m.focusedCol {
			return
		}
	}

	// If focused column is not visible (after the visible columns), scroll right
	if m.focusedCol >= m.scrollOffsetCol+visibleColCount {
		m.scrollOffsetCol = m.focusedCol
	}
}

// halfPageUp moves the focus up half a page.
func (m *Model) halfPageUp() {
	if m.height <= 2 {
		m.moveUp()
		return
	}

	visibleRows := (m.height - 2) / 2 // Half of visible rows
	m.focusedRow -= visibleRows
	if m.focusedRow < 0 {
		m.focusedRow = 0
	}
	m.updateScrollRow()
}

// halfPageDown moves the focus down half a page.
func (m *Model) halfPageDown() {
	if m.height <= 2 {
		m.moveDown()
		return
	}

	visibleRows := (m.height - 2) / 2 // Half of visible rows
	m.focusedRow += visibleRows
	if m.focusedRow >= len(m.rows) {
		m.focusedRow = len(m.rows) - 1
	}
	m.updateScrollRow()
}

// yankCell copies the currently focused cell value to the clipboard.
func (m Model) yankCell() tea.Cmd {
	cellValue := m.SelectedCell()
	if cellValue == "" {
		return nil
	}

	// Escape single quotes in the cell value
	escapedValue := strings.ReplaceAll(cellValue, "'", "'\\''")

	// Try xclip (Linux) first, fall back to pbcopy (macOS), then wl-copy (Wayland)
	cmd := exec.Command("sh", "-c",
		"printf '%s' '"+escapedValue+"' | xclip -selection clipboard 2>/dev/null || "+
			"printf '%s' '"+escapedValue+"' | pbcopy 2>/dev/null || "+
			"printf '%s' '"+escapedValue+"' | wl-copy 2>/dev/null")

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return YankMsg{Success: false, Value: cellValue}
		}
		return YankMsg{Success: true, Value: cellValue}
	})
}

// YankMsg is sent when a cell value is yanked (copied).
type YankMsg struct {
	Success bool
	Value   string
}

// SortChangeMsg is sent when the sort order changes.
type SortChangeMsg struct {
	SortOrders []OrderCol
}

// toggleSort toggles sort for the focused column.
// If column is not sorted: add ASC
// If column is ASC: change to DESC
// If column is DESC: remove from sort
// 's' key adds the column to sort list (multi-column sort)
func (m Model) toggleSort() (Model, tea.Cmd) {
	if len(m.cols) == 0 || m.focusedCol < 0 || m.focusedCol >= len(m.cols) {
		return m, nil
	}

	// Find if this column is already in sort orders
	existingIdx := -1
	for i, sort := range m.orderColumns {
		if sort.ColumnIndex == m.focusedCol {
			existingIdx = i
			break
		}
	}

	if existingIdx == -1 {
		// Column not sorted, add as ASC
		m.orderColumns = append(m.orderColumns, OrderCol{
			ColumnIndex: m.focusedCol,
			Direction:   SortAsc,
		})
	} else if m.orderColumns[existingIdx].Direction == SortAsc {
		// Column is ASC, change to DESC
		m.orderColumns[existingIdx].Direction = SortDesc
	} else {
		// Column is DESC, remove from sort
		m.orderColumns = append(m.orderColumns[:existingIdx], m.orderColumns[existingIdx+1:]...)
	}

	return m, func() tea.Msg {
		return SortChangeMsg{SortOrders: m.orderColumns}
	}
}

func (m *Model) SetSortVisually(sortOrders []OrderCol) {
	m.orderColumns = sortOrders
}

// clearSort clears all sorting.
func (m Model) clearSort() (Model, tea.Cmd) {
	if len(m.orderColumns) == 0 {
		return m, nil
	}

	m.orderColumns = nil

	return m, func() tea.Msg {
		return SortChangeMsg{SortOrders: nil}
	}
}

func (m *Model) ScrollToBottom() {
	if len(m.rows) > 0 {
		m.focusedRow = len(m.rows) - 1
		m.updateScrollRow()
	}
}

func (m *Model) ScrollToTop() {
	m.focusedRow = 0
	m.scrollOffsetRow = 0
}
