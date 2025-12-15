package sharedcomponents

import (
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type TableQuery struct {
	BaseQuery  string
	Limit      int
	SortOrders []OrderByClause
	Offset     int
	SQLResult  *SQLResult
}

func NewTableQuery(baseQuery string, limit int) *TableQuery {
	return &TableQuery{
		BaseQuery: baseQuery,
		Limit:     limit + 1,
		Offset:    0,
	}
}

type ReapplyTableQueryMsg struct {
	Query QueryCompiler
}

func (q *TableQuery) Compile() string {
	// TODO: Rewrite to strings.Builder for efficiency
	fullQuery := q.BaseQuery
	if fullQuery[len(fullQuery)-1] == ';' {
		fullQuery = fullQuery[:len(fullQuery)-1]
	}
	if len(q.SortOrders) > 0 {

		var orderByParts []string
		for _, sort := range q.SortOrders {
			quotedColumn := "\"" + sort.ColumnName + "\""
			orderByParts = append(orderByParts, quotedColumn+" "+sort.Direction)
		}

		orderByClause := " ORDER BY " + strings.Join(orderByParts, ", ")
		fullQuery = fullQuery + orderByClause
	}
	if q.Limit > 0 {
		fullQuery = fmt.Sprintf("%s LIMIT %d", fullQuery, q.Limit)
	}
	if q.Offset > 0 {
		fullQuery = fmt.Sprintf("%s OFFSET %d", fullQuery, q.Offset)
	}
	return fullQuery
}

func (q *TableQuery) HandleSortChange(orderBy []OrderByClause) QueryCompiler {
	log.Printf("Handling sort change: %+v", orderBy)
	q.SortOrders = orderBy
	return q
}

func (q *TableQuery) GetSortOrders() []OrderByClause {
	return q.SortOrders
}

func (q *TableQuery) SetSQLResult(msg *SQLResultMsg) *SQLResult {
	canFetchTotal := len(msg.Rows) > 500
	q.SQLResult = &SQLResult{
		Rows:          msg.Rows,
		Columns:       msg.Columns,
		TotalFetched:  len(msg.Rows) + q.Offset,
		CanFetchTotal: canFetchTotal,
	}
	return q.SQLResult
}

func (q *TableQuery) GetSQLResult() *SQLResult {
	return q.SQLResult
}

// HasNextPage checks if there is a next page available
func (q *TableQuery) HasNextPage() bool {
	return q.SQLResult != nil && len(q.SQLResult.Rows) > q.Limit-1
}

func (q *TableQuery) HasPreviousPage() bool {
	return q.Offset > 0
}

// NextPage increments the offset and returns a message to reapply the query
func (q *TableQuery) NextPage() tea.Cmd {
	if !q.HasNextPage() {
		log.Println("No next page available")
		return nil
	}
	q.Offset += q.Limit - 1
	return func() tea.Msg {
		return ReapplyTableQueryMsg{
			Query: q,
		}
	}
}

func (q *TableQuery) PreviousPage() tea.Cmd {
	if !q.HasPreviousPage() {
		log.Println("No previous page available")
		return nil
	}
	q.Offset -= max(0, q.Limit-1)
	return func() tea.Msg {
		return ReapplyTableQueryMsg{
			Query: q,
		}
	}
}

func (q *TableQuery) Rows() [][]any {
	if len(q.SQLResult.Rows) > q.Limit-1 {
		log.Printf("Returning rows for current page: limit=%d", q.Limit)
		return q.SQLResult.Rows[:q.Limit-1]
	}
	return q.SQLResult.Rows
}

func (q *TableQuery) PageOffset() int {
	return q.Offset
}
