package dbtree

import (
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
)

type handleDBSelectionResult struct {
	err     error
	schemas []*database.Schema
}

type handleSchemaSelectionResult struct {
	tables []*database.Table
	cmd    tea.Cmd
}

type loadTablesColumnsResult struct {
	columns    map[string][]*database.Column
	schemaName string
	databaseID string
	err        error
}
