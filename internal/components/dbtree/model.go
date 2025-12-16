package dbtree

import (
	"github.com/SavingFrame/dbettier/internal/database"
)

type DBTreeModel struct {
	tree     TreeState
	cursor   TreeCursor
	search   TreeSearch
	viewport Viewport
	registry *database.DBRegistry
}

func DBTreeScreen(registry *database.DBRegistry) DBTreeModel {
	var dbNodes []*databaseNode
	for _, db := range registry.GetAll() {
		dbNodes = append(dbNodes, &databaseNode{
			name:     db.Database,
			host:     db.Host,
			expanded: false,
			id:       db.ID,
			db:       db,
		})
	}

	cursor := TreeCursor{path: []int{0}}

	return DBTreeModel{
		cursor:   cursor,
		tree:     TreeState{databases: dbNodes, cursor: &cursor},
		search:   TreeSearch{matchIndex: -1},
		registry: registry,
		viewport: Viewport{},
	}
}

func (m *DBTreeModel) SetSize(width, height int) {
	m.viewport.SetSize(width, height)
}
