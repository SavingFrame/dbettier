package query

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

type BasicSQLQuery struct {
	Query      string
	SortOrders OrderByClauses
	SQLResult  *SQLResult
	// page
	localOffset int
	localRows   [][]any
}

func NewBasicSQLQuery(query string) *BasicSQLQuery {
	return &BasicSQLQuery{
		Query:       query,
		localOffset: 0,
	}
}

func (q *BasicSQLQuery) Compile() string {
	return q.Query
}

func (q *BasicSQLQuery) HandleSortChange(orderBy OrderByClauses) tea.Cmd {
	q.SortOrders = orderBy
	return nil
}

func (q *BasicSQLQuery) GetSortOrders() OrderByClauses {
	return q.SortOrders
}

func (q *BasicSQLQuery) SetSQLResult(msg *SQLResultMsg) *SQLResult {
	q.SQLResult = &SQLResult{
		Rows:          msg.Rows,
		Columns:       msg.Columns,
		TotalFetched:  len(msg.Rows),
		Total:         len(msg.Rows),
		CanFetchTotal: false,
	}
	if len(msg.Rows) > 500 {
		q.localRows = q.SQLResult.Rows[:500]
	} else {
		q.localRows = q.SQLResult.Rows
	}
	return q.SQLResult
}

func (q *BasicSQLQuery) GetSQLResult() *SQLResult {
	return q.SQLResult
}

func (q *BasicSQLQuery) HasNextPage() bool {
	return q.SQLResult != nil && len(q.localRows) >= 500
}

func (q *BasicSQLQuery) HasPreviousPage() bool {
	return q.localOffset > 0
}

func (q *BasicSQLQuery) PreviousPage() tea.Cmd {
	if !q.HasPreviousPage() {
		log.Println("No previous page available")
		return nil
	}
	q.localOffset = max(0, q.localOffset-500)
	return func() tea.Msg {
		return UpdateTableMsg{
			Query: q,
		}
	}
}

func (q *BasicSQLQuery) NextPage() tea.Cmd {
	q.localRows = q.SQLResult.Rows[:500]
	return func() tea.Msg {
		return UpdateTableMsg{
			Query: q,
		}
	}
}

func (q *BasicSQLQuery) Rows() [][]any {
	return q.localRows
}

func (q *BasicSQLQuery) PageOffset() int {
	return q.localOffset
}
