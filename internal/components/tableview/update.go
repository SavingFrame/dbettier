package tableview

import (
	"fmt"
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/pkgs/table"
	"github.com/jackc/pgx/v5/pgtype"
)

// KeyMap defines keybindings for the table view component
type KeyMap struct {
	Enter        key.Binding
	Quit         key.Binding
	NextPage     key.Binding
	PreviousPage key.Binding
}

// DefaultKeyMap returns the default keybindings for the table view
var DefaultKeyMap = KeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select row"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "bottom/next page"),
	),
	PreviousPage: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top/prev page"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextPage, k.PreviousPage, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPage, k.PreviousPage},
		{k.Enter, k.Quit},
	}
}

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case sharedcomponents.SQLResultMsg:
		m.handleSQLResultMsg(msg)
		return m, nil
	case sharedcomponents.UpdateTableMsg:
		m.query = msg.Query
		m.updateTableData(m.query.GetSQLResult())
		return m, nil

	case table.SortChangeMsg:
		// Table sort changed, emit OrderByChangeMsg for SQL layer
		return m, m.handleSortChange(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.nextPageClicked = false
			m.previousPageClicked = false
			m.customMessage = ""
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Quit):
			if m.nextPageClicked || m.previousPageClicked {
				m.nextPageClicked = false
				m.previousPageClicked = false
				m.customMessage = ""
				return m, nil
			}
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.NextPage):
			m.table.ScrollToBottom()
			if m.table.IsLatestRowFocused() && m.query.HasNextPage() {
				switch m.nextPageClicked {
				case true:
					m.nextPageClicked = false
					m.customMessage = ""
					return m, m.query.NextPage()
				case false:
					m.nextPageClicked = true
					m.customMessage = "Click G one more time to go to the next page"
				}
			} else {
				// Clear message when not at bottom or no next page
				m.customMessage = ""
			}
			return m, nil
		case key.Matches(msg, DefaultKeyMap.PreviousPage):
			m.table.ScrollToTop()
			if m.table.IsFirstRowFocused() && m.query.HasPreviousPage() {
				switch m.previousPageClicked {
				case true:
					m.previousPageClicked = false
					m.customMessage = ""
					return m, m.query.PreviousPage()
				case false:
					m.previousPageClicked = true
					m.customMessage = "Click g one more time to go to the previous page"
				}
			}
		default:
			m.nextPageClicked = false
			m.customMessage = ""
		}
	case tea.WindowSizeMsg:
		// Size will be handled by root screen
		return m, nil
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TableViewModel) handleSortChange(msg table.SortChangeMsg) tea.Cmd {
	q := m.query
	if m.databaseID == "" {
		log.Println("Cannot sort: no base query or database ID")
		return nil
	}

	var orderByClauses []sharedcomponents.OrderByClause
	columns := m.table.Columns()

	for _, sort := range msg.SortOrders {
		if sort.ColumnIndex < 0 || sort.ColumnIndex >= len(columns) {
			continue
		}

		orderByClauses = append(orderByClauses, sharedcomponents.OrderByClause{
			ColumnName: columns[sort.ColumnIndex].Title,
			Direction:  sort.Direction.String(),
		})
	}
	log.Printf("TableViewModel: handling sort change: %+v", orderByClauses)
	switch tq := q.(type) {
	case *sharedcomponents.TableQuery:
		tq.HandleSortChange(orderByClauses)
		return func() tea.Msg {
			return sharedcomponents.ReapplyTableQueryMsg{
				Query: tq,
			}
		}
	}
	return nil
}

func formatCellValue(cell any) string {
	switch v := cell.(type) {
	case pgtype.Numeric:
		vTmp, err := v.Value()
		if err != nil {
			return fmt.Sprintf("ERR: %v", err)
		}
		return fmt.Sprintf("%v", vTmp)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (m *TableViewModel) handleSQLResultMsg(msg sharedcomponents.SQLResultMsg) {
	m.query = msg.Query
	m.databaseID = msg.DatabaseID

	sqlResult := msg.Query.SetSQLResult(&msg)
	m.updateTableData(sqlResult)
}

func (m *TableViewModel) updateTableData(r *sharedcomponents.SQLResult) {
	var columns []table.Column
	var rows []table.Row

	// Use a reasonable fixed width for each column (allow scrolling for many columns)
	const minColWidth = 3
	const maxColWidth = 50

	var colSize []int
	for range r.Columns {
		colSize = append(colSize, minColWidth)
	}
	for _, rowData := range m.query.Rows() {
		var rowCells []string
		for cellI, cell := range rowData {
			vText := formatCellValue(cell)
			rowCells = append(rowCells, vText)
			colSize[cellI] = max(len(vText), colSize[cellI])
		}
		rows = append(rows, table.Row(rowCells))
	}

	for colI, colName := range r.Columns {
		// Calculate column width based on column name length, with min/max bounds
		// colWidth := max(min(colSize[colI], maxColWidth), minColWidth) + 2
		colWidth := max(colSize[colI], len(colName)) + 5
		colWidth = min(colWidth, maxColWidth)
		colWidth = max(colWidth, minColWidth)

		columns = append(columns, table.Column{
			Title: colName,
			Width: colWidth,
		})
	}

	m.canFetchTotal = r.CanFetchTotal

	m.table.SetRows(nil)
	m.table.SetColumns(columns)

	m.table.SetRows(rows)
}
