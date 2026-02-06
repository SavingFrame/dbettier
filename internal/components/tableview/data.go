package tableview

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/query"
	"github.com/SavingFrame/dbettier/pkgs/table"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	minColWidth = 3
	maxColWidth = 50
)

// DataState manages query results and table data
type DataState struct {
	query         query.ExecutableQuery
	databaseID    string
	canFetchTotal bool
}

func (d *DataState) Query() query.ExecutableQuery {
	return d.query
}

func (d *DataState) DatabaseID() string {
	return d.databaseID
}

func (d *DataState) CanFetchTotal() bool {
	return d.canFetchTotal
}

func (d *DataState) HasQuery() bool {
	return d.query != nil
}

func (d *DataState) HasNextPage() bool {
	if d.query == nil {
		return false
	}
	return d.query.HasNextPage()
}

func (d *DataState) HasPreviousPage() bool {
	if d.query == nil {
		return false
	}
	return d.query.HasPreviousPage()
}

func (d *DataState) PageOffset() int {
	if d.query == nil {
		return 0
	}
	return d.query.PageOffset()
}

func (d *DataState) GetSortOrders() query.OrderByClauses {
	if d.query == nil {
		return nil
	}
	return d.query.GetSortOrders()
}

func (d *DataState) SetFromSQLResult(msg query.SQLResultMsg) *query.SQLResult {
	d.query = msg.Query
	d.databaseID = msg.DatabaseID
	return msg.Query.SetSQLResult(&msg)
}

func (d *DataState) SetQuery(query query.ExecutableQuery) {
	d.query = query
}

func (d *DataState) BuildTableData(result *query.SQLResult) ([]table.Column, []table.Row) {
	if result == nil || d.query == nil {
		return nil, nil
	}

	// Calculate column sizes based on data
	colSizes := make([]int, len(result.Columns))
	for i := range colSizes {
		colSizes[i] = minColWidth
	}

	var rows []table.Row
	for _, rowData := range d.query.Rows() {
		var rowCells []string
		for cellIdx, cell := range rowData {
			text := formatCellValue(cell)
			rowCells = append(rowCells, text)
			if len(text) > colSizes[cellIdx] {
				colSizes[cellIdx] = len(text)
			}
		}
		rows = append(rows, table.Row(rowCells))
	}

	// Build columns with calculated widths
	var columns []table.Column
	for colIdx, colName := range result.Columns {
		width := max(colSizes[colIdx], len(colName)) + 5
		width = min(width, maxColWidth)
		width = max(width, minColWidth)

		columns = append(columns, table.Column{
			Title: colName,
			Width: width,
		})
	}

	d.canFetchTotal = result.CanFetchTotal

	return columns, rows
}

func formatCellValue(cell any) string {
	switch v := cell.(type) {
	case pgtype.Numeric:
		val, err := v.Value()
		if err != nil {
			return fmt.Sprintf("ERR: %v", err)
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (d *DataState) HandleSortChange(columns []table.Column, sortOrders []table.OrderCol) query.OrderByClauses {
	var orderByClauses query.OrderByClauses

	for _, sort := range sortOrders {
		if sort.ColumnIndex < 0 || sort.ColumnIndex >= len(columns) {
			continue
		}

		orderByClauses = append(orderByClauses, query.OrderByClause{
			ColumnName: columns[sort.ColumnIndex].Title,
			Direction:  sort.Direction.String(),
		})
	}

	return orderByClauses
}

func (d *DataState) RefreshQuery() tea.Cmd {
	return func() tea.Msg {
		return query.ReapplyTableQueryMsg{
			Query: d.query,
		}
	}
}

// IsTableQuery returns true if the current query is a TableQuery (vs BasicQuery)
func (d *DataState) IsTableQuery() bool {
	_, ok := d.query.(*query.TableQuery)
	return ok
}
