package query

type SQLResultMsg struct {
	Rows       [][]any
	Columns    []string // Maybe change, set types for columns, etc
	Query      ExecutableQuery
	DatabaseID string
}

// TODO: I dont know what is it doing
type UpdateTableMsg struct {
	Query ExecutableQuery
}

type ReapplyTableQueryMsg struct {
	Query ExecutableQuery
}
