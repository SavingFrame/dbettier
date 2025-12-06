package sharedcomponents

type BasicSQLQuery struct {
	Query      string
	SortOrders []OrderByClause
	SQLResult  *SQLResult
}

func (q *BasicSQLQuery) Compile() string {
	return q.Query
}

func (q *BasicSQLQuery) HandleSortChange(orderBy []OrderByClause) QueryCompiler {
	q.SortOrders = orderBy
	return q
}

func (q *BasicSQLQuery) GetSortOrders() []OrderByClause {
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
	return q.SQLResult
}

func (q *BasicSQLQuery) HasNextPage() bool {
	return q.SQLResult != nil && q.SQLResult.Total >= 500
}
