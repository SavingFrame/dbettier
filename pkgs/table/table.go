// Package table provides a table component with cell-level focus (row + column)
// for Bubble Tea applications.
package table

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// SortDirection represents the direction of sorting.
type SortDirection int

const (
	SortAsc SortDirection = iota
	SortDesc
)

// String returns the string representation of SortDirection.
func (s SortDirection) String() string {
	if s == SortAsc {
		return "ASC"
	}
	return "DESC"
}

// OrderCol represents a column sort order.
type OrderCol struct {
	ColumnIndex int
	Direction   SortDirection
}

// Model defines the state for the table widget with cell-level focus.
type Model struct {
	cols []Column
	rows []Row

	// Cell-level focus
	focusedRow int
	focusedCol int

	// Scroll offsets
	scrollOffsetRow int
	scrollOffsetCol int

	// Dimensions
	width  int
	height int

	// State
	focused bool
	styles  Styles

	// Sorting
	orderColumns []OrderCol
}

// Row represents one line in the table.
type Row []string

// Column defines the table structure.
type Column struct {
	Title string
	Width int
}

// Styles contains style definitions for the table component.
type Styles struct {
	Header       lipgloss.Style
	Cell         lipgloss.Style
	SelectedCell lipgloss.Style
	SelectedRow  lipgloss.Style
	SelectedCol  lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("240")),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		SelectedCell: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("229")).
			Bold(true),
		SelectedRow: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("237")),
		SelectedCol: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("237")),
	}
}

// Option is used to set options in New.
type Option func(*Model)

// New creates a new model for the table widget.
func New(opts ...Option) Model {
	m := Model{
		focusedRow:      0,
		focusedCol:      0,
		scrollOffsetRow: 0,
		scrollOffsetCol: 0,
		focused:         false,
		styles:          DefaultStyles(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithColumns sets the table columns (headers).
func WithColumns(cols []Column) Option {
	return func(m *Model) {
		m.cols = cols
	}
}

// WithRows sets the table rows (data).
func WithRows(rows []Row) Option {
	return func(m *Model) {
		m.rows = rows
	}
}

// WithHeight sets the height of the table.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.height = h
	}
}

// WithWidth sets the width of the table.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// WithFocused sets the focus state of the table.
func WithFocused(f bool) Option {
	return func(m *Model) {
		m.focused = f
	}
}

// WithStyles sets the table styles.
func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

// SetColumns updates the table columns.
func (m *Model) SetColumns(cols []Column) {
	m.cols = cols
	// Reset column focus if out of bounds
	if m.focusedCol >= len(m.cols) {
		m.focusedCol = 0
	}
}

// SetRows updates the table rows.
func (m *Model) SetRows(rows []Row) {
	m.rows = rows
	// Reset row focus if out of bounds
	if m.focusedRow >= len(m.rows) {
		m.focusedRow = 0
	}
}

// SetHeight updates the table height.
func (m *Model) SetHeight(h int) {
	m.height = h
}

// SetWidth updates the table width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetStyles updates the table styles.
func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

// Focused returns the current focus state.
func (m Model) Focused() bool {
	return m.focused
}

// Focus sets the focus state to true.
func (m *Model) Focus() {
	m.focused = true
}

// Blur sets the focus state to false.
func (m *Model) Blur() {
	m.focused = false
}

// Columns returns the current columns.
func (m Model) Columns() []Column {
	return m.cols
}

// Rows returns the current rows.
func (m Model) Rows() []Row {
	return m.rows
}

// SelectedRow returns the currently focused row data.
func (m Model) SelectedRow() Row {
	if m.focusedRow >= 0 && m.focusedRow < len(m.rows) {
		return m.rows[m.focusedRow]
	}
	return nil
}

// SelectedCell returns the currently focused cell data.
func (m Model) SelectedCell() string {
	if m.focusedRow >= 0 && m.focusedRow < len(m.rows) {
		row := m.rows[m.focusedRow]
		if m.focusedCol >= 0 && m.focusedCol < len(row) {
			return row[m.focusedCol]
		}
	}
	return ""
}

// FocusedPosition returns the current focused row and column indices.
func (m Model) FocusedPosition() (row, col int) {
	return m.focusedRow, m.focusedCol
}

// OrderColumns returns the current sort orders.
func (m Model) OrderColumns() []OrderCol {
	return m.orderColumns
}

// Init initializes the table.
func (m Model) Init() tea.Cmd {
	return nil
}
