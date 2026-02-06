package query

import tea "charm.land/bubbletea/v2"

type ExecutableQuery interface {
	Compile() string
	HandleSortChange(orderBy OrderByClauses) tea.Cmd
	GetSortOrders() OrderByClauses
	SetSQLResult(*SQLResultMsg) *SQLResult
	GetSQLResult() *SQLResult
	HasPreviousPage() bool
	HasNextPage() bool
	NextPage() tea.Cmd
	PreviousPage() tea.Cmd
	Rows() [][]any
	PageOffset() int
}
