package messages

import (
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
)

// OpenQueryTabMsg creates new basic query tab for database. You can just open empty tab if you pass QueryCompiler with empty query
type OpenQueryTabMsg struct {
	Query      sharedcomponents.QueryCompiler
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
