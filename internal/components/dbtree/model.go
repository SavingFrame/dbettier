package dbtree

import (
	"database/sql"

	"github.com/SavingFrame/dbettier/internal/database"
)

type tableColumnNode struct {
	name      string
	dataType  string
	maxLength sql.NullInt32
}

type schemaTableNode struct {
	name     string
	columns  []*tableColumnNode
	expanded bool
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
	id       string
	parsed   bool
}

// TreeLevel represents the depth in the tree hierarchy
type TreeLevel int

const (
	DatabaseLevel TreeLevel = iota
	SchemaLevel
	TableLevel
	TableColumnLevel
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
	windowWidth  int
	scrollOffset int
}

func DBTreeScreen(registry *database.DBRegistry) DBTreeModel {
	var dbNodes []*databaseNode
	for _, db := range registry.GetAll() {
		dbNodes = append(dbNodes, &databaseNode{
			name:     db.Database,
			host:     db.Host,
			expanded: false,
			id:       db.ID,
		})
	}

	return DBTreeModel{
		cursor: treeCursor{
			path: []int{0}, // Start at first database
		},
		databases:    dbNodes,
		registry:     registry,
		windowHeight: 0,
		windowWidth:  0,
		scrollOffset: 0,
	}
}

// SetSize updates the dimensions of the DBTree view
func (m *DBTreeModel) SetSize(width, height int) {
	m.windowWidth = width
	m.windowHeight = height
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

func (c treeCursor) tableColumnIndex() int {
	if len(c.path) > 3 {
		return c.path[3]
	}
	return -1
}

// TODO: Remove
// isAtDatabaseLevel returns true if cursor is on a database (not a schema)
func (c treeCursor) isAtDatabaseLevel() bool {
	return c.atLevel(DatabaseLevel)
}

// Node navigation helper methods

// getDatabase safely returns a database node by index
func (m DBTreeModel) getDatabase(dbIdx int) *databaseNode {
	if dbIdx < 0 || dbIdx >= len(m.databases) {
		return nil
	}
	return m.databases[dbIdx]
}

// getCurrentDatabase returns the database at the current cursor position
func (m DBTreeModel) getCurrentDatabase() *databaseNode {
	return m.getDatabase(m.cursor.dbIndex())
}

// getSchema safely returns a schema node by indices
func (m DBTreeModel) getSchema(dbIdx, schemaIdx int) *databaseSchemaNode {
	db := m.getDatabase(dbIdx)
	if db == nil || schemaIdx < 0 || schemaIdx >= len(db.schemas) {
		return nil
	}
	return db.schemas[schemaIdx]
}

// getCurrentSchema returns the schema at the current cursor position
func (m DBTreeModel) getCurrentSchema() *databaseSchemaNode {
	if m.cursor.level() < SchemaLevel {
		return nil
	}
	return m.getSchema(m.cursor.dbIndex(), m.cursor.schemaIndex())
}

// getTable safely returns a table node by indices
func (m DBTreeModel) getTable(dbIdx, schemaIdx, tableIdx int) *schemaTableNode {
	schema := m.getSchema(dbIdx, schemaIdx)
	if schema == nil || tableIdx < 0 || tableIdx >= len(schema.tables) {
		return nil
	}
	return schema.tables[tableIdx]
}

// getCurrentTable returns the table at the current cursor position
func (m DBTreeModel) getCurrentTable() *schemaTableNode {
	if m.cursor.level() < TableLevel {
		return nil
	}
	return m.getTable(m.cursor.dbIndex(), m.cursor.schemaIndex(), m.cursor.tableIndex())
}

// getColumn safely returns a column node by indices
func (m DBTreeModel) getColumn(dbIdx, schemaIdx, tableIdx, colIdx int) *tableColumnNode {
	table := m.getTable(dbIdx, schemaIdx, tableIdx)
	if table == nil || colIdx < 0 || colIdx >= len(table.columns) {
		return nil
	}
	return table.columns[colIdx]
}

// getCurrentColumn returns the column at the current cursor position
func (m DBTreeModel) getCurrentColumn() *tableColumnNode {
	if m.cursor.level() < TableColumnLevel {
		return nil
	}
	return m.getColumn(m.cursor.dbIndex(), m.cursor.schemaIndex(), m.cursor.tableIndex(), m.cursor.tableColumnIndex())
}

// findDatabase finds a database by ID and returns it with its index
func (m DBTreeModel) findDatabase(id string) (*databaseNode, int) {
	for i, db := range m.databases {
		if db.id == id {
			return db, i
		}
	}
	return nil, -1
}

// findSchema finds a schema by name within a database
func (m DBTreeModel) findSchema(db *databaseNode, name string) (*databaseSchemaNode, int) {
	if db == nil {
		return nil, -1
	}
	for i, schema := range db.schemas {
		if schema.name == name {
			return schema, i
		}
	}
	return nil, -1
}

// hasChildren returns true if the current cursor position has expandable children
func (m DBTreeModel) hasChildren() bool {
	switch m.cursor.level() {
	case DatabaseLevel:
		db := m.getCurrentDatabase()
		return db != nil && len(db.schemas) > 0
	case SchemaLevel:
		schema := m.getCurrentSchema()
		return schema != nil && len(schema.tables) > 0
	case TableLevel:
		table := m.getCurrentTable()
		return table != nil && len(table.columns) > 0
	case TableColumnLevel:
		return false
	}
	return false
}

// isExpanded returns true if the current node is expanded
func (m DBTreeModel) isExpanded() bool {
	switch m.cursor.level() {
	case DatabaseLevel:
		db := m.getCurrentDatabase()
		return db != nil && db.expanded
	case SchemaLevel:
		schema := m.getCurrentSchema()
		return schema != nil && schema.expanded
	case TableLevel:
		table := m.getCurrentTable()
		return table != nil && table.expanded
	case TableColumnLevel:
		return false
	}
	return false
}

// getCurrentIndex returns the current index at the current level
func (m DBTreeModel) getCurrentIndex() int {
	switch m.cursor.level() {
	case DatabaseLevel:
		return m.cursor.dbIndex()
	case SchemaLevel:
		return m.cursor.schemaIndex()
	case TableLevel:
		return m.cursor.tableIndex()
	case TableColumnLevel:
		return m.cursor.tableColumnIndex()
	}
	return 0
}

// getLastVisibleDescendantAtPath returns the path to the deepest visible descendant of a node
func (m DBTreeModel) getLastVisibleDescendantAtPath(path []int) []int {
	if len(path) == 0 {
		return path
	}

	// For database level
	if len(path) == 1 {
		db := m.getDatabase(path[0])
		if db == nil || !db.expanded || len(db.schemas) == 0 {
			return path
		}
		// Recurse into last schema
		return m.getLastVisibleDescendantAtPath([]int{path[0], len(db.schemas) - 1})
	}

	// For schema level
	if len(path) == 2 {
		schema := m.getSchema(path[0], path[1])
		if schema == nil || !schema.expanded || len(schema.tables) == 0 {
			return path
		}
		// Recurse into last table
		return m.getLastVisibleDescendantAtPath([]int{path[0], path[1], len(schema.tables) - 1})
	}

	// For table level
	if len(path) == 3 {
		table := m.getTable(path[0], path[1], path[2])
		if table == nil || !table.expanded || len(table.columns) == 0 {
			return path
		}
		// Go to last column
		return []int{path[0], path[1], path[2], len(table.columns) - 1}
	}

	// Column level - no children
	return path
}

// TODO: Refactor to avoid code duplication with rendering logic
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
					for tableIdx, table := range schema.tables {
						if m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.tableIndex() == tableIdx {
							return lineNum
						}
						lineNum++

						if table.expanded && len(table.columns) > 0 {
							// If table is expanded, count columns
							for columnIdx := range table.columns {
								if m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.tableIndex() == tableIdx && m.cursor.tableColumnIndex() == columnIdx {
									return lineNum
								}
								lineNum++
							}
						}
					}
				}
			}
		}
	}
	return lineNum
}
