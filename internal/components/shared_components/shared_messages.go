package sharedcomponents

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SQLQuery struct {
	BaseQuery  string
	Limit      int
	SortOrders []OrderByClause
}

func (q SQLQuery) Compile() string {
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

		// orderByClause := " ORDER BY " + joinStrings(orderByParts, ", ")
		orderByClause := " ORDER BY " + strings.Join(orderByParts, ", ")
		fullQuery = fullQuery + orderByClause
	}
	if q.Limit > 0 {
		fullQuery = fmt.Sprintf("%s LIMIT %d", fullQuery, q.Limit)
	}
	return fullQuery
}

type SetSQLTextMsg struct {
	Query      SQLQuery
	DatabaseID string
}

type SQLResultMsg struct {
	Columns    []string
	Rows       [][]any
	Query      SQLQuery
	DatabaseID string
}

type OrderByClause struct {
	ColumnName string
	Direction  string // "ASC" or "DESC"
}

type ComponentTarget int

const (
	TargetSQLCommandBar ComponentTarget = 1 << iota
	TargetTableView
	TargetDBTree
)

var MessageRoutes = map[string]ComponentTarget{
	"sharedcomponents.SetSQLTextMsg":    TargetSQLCommandBar | TargetTableView,
	"sharedcomponents.SQLResultMsg":     TargetTableView,
	"sharedcomponents.OrderByChangeMsg": TargetSQLCommandBar,
}

func GetMessageType(msg tea.Msg) string {
	return fmt.Sprintf("%T", msg)
}
