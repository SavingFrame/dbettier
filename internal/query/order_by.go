package query

import (
	"fmt"
	"strings"
)

type OrderByClause struct {
	ColumnName string
	Direction  string // "ASC" or "DESC"
}
type OrderByClauses []OrderByClause

func (o OrderByClauses) String() string {
	strs := make([]string, len(o))
	for i, clause := range o {
		strs[i] = fmt.Sprintf("\"%s\" %s", clause.ColumnName, clause.Direction)
	}
	return strings.Join(strs, ", ")
}

func ParseOrderByClauses(s string) (OrderByClauses, error) {
	var clauses OrderByClauses
	for orderClause := range strings.SplitSeq(s, ",") {
		orderClause = strings.TrimSpace(orderClause)
		parts := strings.SplitN(orderClause, " ", 2)
		var columnName, direction string
		switch len(parts) {
		case 2:
			columnName = strings.Trim(parts[0], "\"")
			direction = strings.ToUpper(strings.TrimSpace(parts[1]))
		case 1:
			columnName = strings.Trim(parts[0], "\"")
			direction = "ASC"
		default:
			return nil, fmt.Errorf("invalid order by clause: %s", orderClause)
		}
		clauses = append(clauses, OrderByClause{
			ColumnName: columnName,
			Direction:  direction,
		})
	}
	return clauses, nil
}
