package dbtree

import "github.com/SavingFrame/dbettier/internal/database"

type schemaTableNode struct {
	name string
}

type databaseSchemaNode struct {
	name     string
	tables   []*schemaTableNode
	expanded bool
}

type databaseNode struct {
	name     string
	host     string
	schemas  []*databaseSchemaNode
	expanded bool
}

// TreeLevel represents the depth in the tree hierarchy
type TreeLevel int

const (
	DatabaseLevel TreeLevel = iota
	SchemaLevel
	TableLevel
)

// treeCursor represents the current focus position in the tree using a path
// path[0] = database index
// path[1] = schema index (if at schema level or deeper)
// path[2] = table index (if at table level or deeper)
type treeCursor struct {
	path []int
}

type DBTreeModel struct {
	cursor       treeCursor
	databases    []*databaseNode
	registry     *database.DBRegistry
	windowHeight int
	scrollOffset int
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
			path: []int{0}, // Start at first database
		},
		databases:    dbNodes,
		registry:     registry,
		windowHeight: 20,
		scrollOffset: 0,
	}
}

// level returns the current tree level
func (c treeCursor) level() TreeLevel {
	return TreeLevel(len(c.path) - 1)
}

// atLevel checks if cursor is at a specific level
func (c treeCursor) atLevel(level TreeLevel) bool {
	return c.level() == level
}

// dbIndex returns the database index
func (c treeCursor) dbIndex() int {
	if len(c.path) > 0 {
		return c.path[0]
	}
	return 0
}

// schemaIndex returns the schema index, or -1 if not at schema level or deeper
func (c treeCursor) schemaIndex() int {
	if len(c.path) > 1 {
		return c.path[1]
	}
	return -1
}

// tableIndex returns the table index, or -1 if not at table level
func (c treeCursor) tableIndex() int {
	if len(c.path) > 2 {
		return c.path[2]
	}
	return -1
}

// isAtDatabaseLevel returns true if cursor is on a database (not a schema)
func (c treeCursor) isAtDatabaseLevel() bool {
	return c.atLevel(DatabaseLevel)
}

func (m DBTreeModel) getCursorVisualLine() int {
	lineNum := 1

	for dbIdx, db := range m.databases {
		// Current database is at lineNum
		if m.cursor.dbIndex() == dbIdx && m.cursor.isAtDatabaseLevel() {
			return lineNum
		}
		lineNum++

		// If database is expanded, count schemas
		if db.expanded && len(db.schemas) > 0 {
			for schemaIdx, schema := range db.schemas {
				if m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.atLevel(SchemaLevel) {
					return lineNum
				}
				lineNum++

				// If schema is expanded, count tables
				if schema.expanded && len(schema.tables) > 0 {
					for tableIdx := range schema.tables {
						if m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.tableIndex() == tableIdx {
							return lineNum
						}
						lineNum++
					}
				}
			}
		}
	}
	return lineNum
}
