package tableview

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/pkgs/table"
)

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case sharedcomponents.SQLResultMsg:
		m.handleSQLResultMsg(msg)
		return m, nil

	case table.SortChangeMsg:
		// Table sort changed, emit OrderByChangeMsg for SQL layer
		return m, m.handleSortChange(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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

func (m *TableViewModel) handleSQLResultMsg(msg sharedcomponents.SQLResultMsg) {
	m.query = msg.Query
	m.databaseID = msg.DatabaseID

	sqlResult := msg.Query.SetSQLResult(&msg)
	var columns []table.Column
	var rows []table.Row

	// Use a reasonable fixed width for each column (allow scrolling for many columns)
	const minColWidth = 3
	const maxColWidth = 50

	var colSize []int
	for range msg.Columns {
		colSize = append(colSize, minColWidth)
	}
	var rawRows [][]any
	if len(sqlResult.Rows) > 500 {
		rawRows = sqlResult.Rows[:500]
	} else {
		rawRows = sqlResult.Rows
	}
	for _, rowData := range rawRows {
		var rowCells []string
		for cellI, cell := range rowData {
			v := fmt.Sprintf("%v", cell)
			rowCells = append(rowCells, v)
			colSize[cellI] = max(len(v), colSize[cellI])
		}
		rows = append(rows, table.Row(rowCells))
	}

	for colI, colName := range msg.Columns {
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

	m.canFetchTotal = sqlResult.CanFetchTotal

	m.table.SetRows(nil)
	m.table.SetColumns(columns)

	m.table.SetRows(rows)
}
