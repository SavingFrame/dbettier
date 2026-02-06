package tableview

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/query"
	"github.com/SavingFrame/dbettier/internal/theme"
	zone "github.com/lrstanley/bubblezone/v2"
)

// inputStyles returns the text input styles using current theme
func inputStyles() textinput.Styles {
	colors := theme.Current().Colors
	inputBg := colors.Surface
	inputTextStyle := lipgloss.NewStyle().Background(inputBg).Foreground(colors.Text)
	inputPlaceholderStyle := lipgloss.NewStyle().Background(inputBg).Foreground(colors.Muted)
	inputPromptStyle := lipgloss.NewStyle().Background(inputBg).Foreground(colors.Info)

	return textinput.Styles{
		Focused: textinput.StyleState{
			Text:        inputTextStyle,
			Placeholder: inputPlaceholderStyle,
			Prompt:      inputPromptStyle,
		},
		Blurred: textinput.StyleState{
			Text:        inputTextStyle,
			Placeholder: inputPlaceholderStyle,
			Prompt:      inputPromptStyle.Foreground(colors.Muted),
		},
	}
}

// StatusBarFocus represents which element of the status bar is focused
type StatusBarFocus int

const (
	StatusBarFocusNone StatusBarFocus = iota
	StatusBarFocusFilter
	StatusBarFocusOrdering
)

// StatusBar handles the status bar UI for the table view
type StatusBar struct {
	pagination Pagination

	// Display state (synced from tableview)
	focusedRow   int
	totalRows    int
	pageOffset   int
	canFetchMore bool
	focusedCol   int
	totalCols    int
	isTableQuery bool
	sortOrders   query.OrderByClauses

	// Input fields
	filterInput   textinput.Model
	orderingInput textinput.Model
	focus         StatusBarFocus

	// Dimensions
	width int
}

// NewStatusBar creates a new status bar
func NewStatusBar() StatusBar {
	styles := inputStyles()

	filterInput := textinput.New()
	filterInput.Placeholder = "filter..."
	filterInput.CharLimit = 256
	filterInput.SetWidth(25)
	filterInput.SetStyles(styles)

	orderingInput := textinput.New()
	orderingInput.Placeholder = "col ASC/DESC"
	orderingInput.CharLimit = 256
	orderingInput.SetWidth(25)
	orderingInput.SetStyles(styles)

	return StatusBar{
		pagination:    Pagination{},
		filterInput:   filterInput,
		orderingInput: orderingInput,
		focus:         StatusBarFocusNone,
	}
}

// Pagination returns a pointer to the pagination state
func (s *StatusBar) Pagination() *Pagination {
	return &s.pagination
}

// SetWidth sets the width of the status bar
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// Focus returns the current focus state
func (s *StatusBar) Focus() StatusBarFocus {
	return s.focus
}

// SetFocus sets the focus to a specific element
func (s *StatusBar) SetFocus(focus StatusBarFocus) {
	s.focus = focus
	switch focus {
	case StatusBarFocusFilter:
		s.filterInput.Focus()
		s.orderingInput.Blur()
	case StatusBarFocusOrdering:
		s.orderingInput.Focus()
		s.filterInput.Blur()
	default:
		s.filterInput.Blur()
		s.orderingInput.Blur()
	}
}

// IsFocused returns true if any input in the status bar is focused
func (s *StatusBar) IsFocused() bool {
	return s.focus != StatusBarFocusNone
}

// FilterInput returns a pointer to the filter input
func (s *StatusBar) FilterInput() *textinput.Model {
	return &s.filterInput
}

// OrderingInput returns a pointer to the ordering input
func (s *StatusBar) OrderingInput() *textinput.Model {
	return &s.orderingInput
}

// FilterValue returns the current filter value
func (s *StatusBar) FilterValue() string {
	return s.filterInput.Value()
}

// OrderingValue returns the current ordering value
func (s *StatusBar) OrderingValue() string {
	return s.orderingInput.Value()
}

// SyncState updates the status bar display state from tableview
func (s *StatusBar) SyncState(
	focusedRow, totalRows, pageOffset int,
	canFetchMore bool,
	focusedCol, totalCols int,
	isTableQuery bool,
	sortOrders query.OrderByClauses,
) {
	s.focusedRow = focusedRow
	s.totalRows = totalRows
	s.pageOffset = pageOffset
	s.canFetchMore = canFetchMore
	s.focusedCol = focusedCol
	s.totalCols = totalCols
	s.isTableQuery = isTableQuery
	s.sortOrders = sortOrders
	s.orderingInput.SetValue(sortOrders.String())
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Single line: Controls on left, position info on right
	controls := s.renderControls()
	posInfo := s.renderPositionInfo()

	// Add pagination message if present
	paginationMsg := ""
	if msg := s.pagination.Message(); msg != "" {
		paginationMsg = "   " + sbPaginationMsgStyle().Render(" "+msg)
	}

	// Calculate spacing to push position info to the right
	leftContent := controls + paginationMsg
	spacing := ""
	if s.width > 0 {
		leftLen := lipgloss.Width(leftContent)
		rightLen := lipgloss.Width(posInfo)
		spaceNeeded := s.width - leftLen - rightLen - 2
		if spaceNeeded > 0 {
			spacing = strings.Repeat(" ", spaceNeeded)
		}
	}

	return leftContent + spacing + posInfo
}

func (s *StatusBar) renderControls() string {
	var parts []string

	// Filter input
	filterLabel := sbInputLabelStyle().Render(" Filter ")
	filterInput := zone.Mark("filterInput", sbInputStyle().Render(s.filterInput.View()))
	parts = append(parts, filterLabel+filterInput)

	// Order input
	orderLabel := sbInputLabelStyle().Render(" Order ")
	orderInput := zone.Mark("orderingInput", sbInputStyle().Render(s.orderingInput.View()))
	parts = append(parts, orderLabel+orderInput)

	// Buttons
	refreshBtn := zone.Mark("refresh", sbButtonPrimaryStyle().Render("↻ Refresh"))
	countBtn := zone.Mark("count", sbButtonStyle().Render("# Count"))

	parts = append(parts, refreshBtn)
	parts = append(parts, countBtn)

	return strings.Join(parts, " ")
}

func (s *StatusBar) renderPositionInfo() string {
	var parts []string

	// Table type icon
	icon := ""
	if s.isTableQuery {
		icon = "󰓫"
	}
	parts = append(parts, sbIconStyle().Render(icon))

	// Sort orders (before position)
	if len(s.sortOrders) > 0 {
		var orderParts []string
		for _, order := range s.sortOrders {
			dir := "↑"
			if order.Direction == "DESC" {
				dir = "↓"
			}
			orderParts = append(orderParts, sbSortStyle().Render(order.ColumnName+" "+dir))
		}
		sortInfo := sbLabelStyle().Render("Sort ") + strings.Join(orderParts, sbDimStyle().Render(", "))
		parts = append(parts, sortInfo)
	}

	// Row position
	if s.totalRows > 0 {
		currentPos := s.focusedRow + 1 + s.pageOffset
		totalRowsStr := fmt.Sprintf("%d", s.totalRows+s.pageOffset)
		if s.canFetchMore {
			totalRowsStr += "+"
		}
		rowInfo := sbLabelStyle().Render("Row ") +
			sbValueStyle().Render(fmt.Sprintf("%d", currentPos)) +
			sbDimStyle().Render("/") +
			sbValueStyle().Render(totalRowsStr)
		parts = append(parts, rowInfo)
	}

	// Col position
	if s.totalCols > 1 {
		colInfo := sbLabelStyle().Render("Col ") +
			sbValueStyle().Render(fmt.Sprintf("%d", s.focusedCol+1)) +
			sbDimStyle().Render("/") +
			sbValueStyle().Render(fmt.Sprintf("%d", s.totalCols))
		parts = append(parts, colInfo)
	}

	sep := sbSepStyle().Render(" │ ")
	return strings.Join(parts, sep)
}
