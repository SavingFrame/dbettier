package dbtree

import "github.com/SavingFrame/dbettier/internal/database"

type databaseSchemaNode struct {
	name string
}

type databaseNode struct {
	name     string
	host     string
	schemas  []*databaseSchemaNode
	expanded bool
}

// treeCursor represents the current focus position in the tree
type treeCursor struct {
	dbIndex     int
	schemaIndex int // -1 if focused on database level
}

type DBTreeModel struct {
	cursor    treeCursor
	databases []*databaseNode
	registry  *database.DBRegistry
}

func DBTreeScreen(registry *database.DBRegistry) DBTreeModel {
	var dbNodes []*databaseNode
	for _, db := range registry.GetAll() {
		dbNodes = append(dbNodes, &databaseNode{
			name:     db.Database,
			host:     db.Host,
			expanded: false,
		})
	}
	return DBTreeModel{
		cursor: treeCursor{
			dbIndex:     0,
			schemaIndex: -1, // Start at database level
		},
		databases: dbNodes,
		registry:  registry,
	}
}

// isAtDatabaseLevel returns true if cursor is on a database (not a schema)
func (c treeCursor) isAtDatabaseLevel() bool {
	return c.schemaIndex == -1
}
