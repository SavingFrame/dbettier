package sharedcomponents

import (
	"fmt"
	"log"
	"strings"
)

type TableQuery struct {
	BaseQuery  string
	Limit      int
	SortOrders []OrderByClause
	Offset     int
	DatabaseID string
	SQLResult  *SQLResult
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
		TotalFetched:  len(msg.Rows),
		CanFetchTotal: canFetchTotal,
	}
	return q.SQLResult
}

func (q *TableQuery) HasNextPage() bool {
	return q.SQLResult == nil && len(q.SQLResult.Rows) > q.Limit-1
}
