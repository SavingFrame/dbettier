package messages

import (
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/SavingFrame/dbettier/internal/query"
)

// OpenQueryTabMsg creates new basic query tab for database. You can just open empty tab if you pass QueryCompiler with empty query
type OpenQueryTabMsg struct {
	Query      query.ExecutableQuery
	DatabaseID string
}

// ExecuteSQLTextMsg executes raw SQL text on database in the current tab.
type ExecuteSQLTextMsg struct {
	Query      string
	DatabaseID string
}

// OpenTableAndExecuteMsg creates a new table tab and executes a query on it
type OpenTableAndExecuteMsg struct {
	Table      *database.Table
	DatabaseID string
}
