package dbtree

import (
	"database/sql"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
)

// TreeLevel represents the depth in the tree hierarchy
type TreeLevel int

type TreeState struct {
	databases []*databaseNode
	cursor    *TreeCursor // Pointer to cursor - avoids passing it everywhere
}

const (
	DatabaseLevel TreeLevel = iota
	SchemaLevel
	TableLevel
	TableColumnLevel
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
	table    *database.Table
}

type databaseSchemaNode struct {
	name     string
	tables   []*schemaTableNode
	expanded bool
	schema   *database.Schema
}

type databaseNode struct {
	name     string
	host     string
	schemas  []*databaseSchemaNode
	expanded bool
	id       string
	parsed   bool
	db       *database.Database
}

// GetDatabase safely returns a database node by index
func (t *TreeState) GetDatabase(dbIdx int) *databaseNode {
	if dbIdx < 0 || dbIdx >= len(t.databases) {
		return nil
	}
	return t.databases[dbIdx]
}

// CurrentDatabase returns the database at the current cursor position
func (t *TreeState) CurrentDatabase() *databaseNode {
	return t.GetDatabase(t.cursor.DbIndex())
}

// GetSchema safely returns a schema node by indices
func (t *TreeState) GetSchema(dbIdx, schemaIdx int) *databaseSchemaNode {
	db := t.GetDatabase(dbIdx)
	if db == nil || schemaIdx < 0 || schemaIdx >= len(db.schemas) {
		return nil
	}
	return db.schemas[schemaIdx]
}

// CurrentSchema returns the schema at the current cursor position
func (t *TreeState) CurrentSchema() *databaseSchemaNode {
	if t.cursor.Level() < SchemaLevel {
		return nil
	}
	return t.GetSchema(t.cursor.DbIndex(), t.cursor.SchemaIndex())
}

// GetTable safely returns a table node by indices
func (t *TreeState) GetTable(dbIdx, schemaIdx, tableIdx int) *schemaTableNode {
	schema := t.GetSchema(dbIdx, schemaIdx)
	if schema == nil || tableIdx < 0 || tableIdx >= len(schema.tables) {
		return nil
	}
	return schema.tables[tableIdx]
}

// CurrentTable returns the table at the current cursor position
func (t *TreeState) CurrentTable() *schemaTableNode {
	if t.cursor.Level() < TableLevel {
		return nil
	}
	return t.GetTable(t.cursor.DbIndex(), t.cursor.SchemaIndex(), t.cursor.TableIndex())
}

// GetColumn safely returns a column node by indices
func (t *TreeState) GetColumn(dbIdx, schemaIdx, tableIdx, colIdx int) *tableColumnNode {
	table := t.GetTable(dbIdx, schemaIdx, tableIdx)
	if table == nil || colIdx < 0 || colIdx >= len(table.columns) {
		return nil
	}
	return table.columns[colIdx]
}

// CurrentColumn returns the column at the current cursor position
func (t *TreeState) CurrentColumn() *tableColumnNode {
	if t.cursor.Level() < TableColumnLevel {
		return nil
	}
	return t.GetColumn(t.cursor.DbIndex(), t.cursor.SchemaIndex(), t.cursor.TableIndex(), t.cursor.TableColumnIndex())
}

// FindDatabase finds a database by ID and returns it with its index
func (t *TreeState) FindDatabase(id string) (*databaseNode, int) {
	for i, db := range t.databases {
		if db.id == id {
			return db, i
		}
	}
	return nil, -1
}

// FindSchema finds a schema by name within a database
func (t *TreeState) FindSchema(db *databaseNode, name string) (*databaseSchemaNode, int) {
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

// HasChildren returns true if the current cursor position has expandable children
func (t *TreeState) HasChildren() bool {
	switch t.cursor.Level() {
	case DatabaseLevel:
		db := t.CurrentDatabase()
		return db != nil && len(db.schemas) > 0
	case SchemaLevel:
		schema := t.CurrentSchema()
		return schema != nil && len(schema.tables) > 0
	case TableLevel:
		table := t.CurrentTable()
		return table != nil && len(table.columns) > 0
	case TableColumnLevel:
		return false
	}
	return false
}

// IsExpanded returns true if the current node is expanded
func (t *TreeState) IsExpanded() bool {
	switch t.cursor.Level() {
	case DatabaseLevel:
		db := t.CurrentDatabase()
		return db != nil && db.expanded
	case SchemaLevel:
		schema := t.CurrentSchema()
		return schema != nil && schema.expanded
	case TableLevel:
		table := t.CurrentTable()
		return table != nil && table.expanded
	case TableColumnLevel:
		return false
	}
	return false
}

// TODO: What is cursor doing here
// getCurrentIndex returns the current index at the current level
func (c *TreeCursor) CurrentIndex() int {
	switch c.Level() {
	case DatabaseLevel:
		return c.DbIndex()
	case SchemaLevel:
		return c.SchemaIndex()
	case TableLevel:
		return c.TableIndex()
	case TableColumnLevel:
		return c.TableColumnIndex()
	}
	return 0
}

// LastVisibleDescendant returns the path to the deepest visible descendant of a node
func (t *TreeState) LastVisibleDescendant(path []int) []int {
	if len(path) == 0 {
		return path
	}

	// For database level
	if len(path) == 1 {
		db := t.GetDatabase(path[0])
		if db == nil || !db.expanded || len(db.schemas) == 0 {
			return path
		}
		// Recurse into last schema
		return t.LastVisibleDescendant([]int{path[0], len(db.schemas) - 1})
	}

	// For schema level
	if len(path) == 2 {
		schema := t.GetSchema(path[0], path[1])
		if schema == nil || !schema.expanded || len(schema.tables) == 0 {
			return path
		}
		// Recurse into last table
		return t.LastVisibleDescendant([]int{path[0], path[1], len(schema.tables) - 1})
	}

	// For table level
	if len(path) == 3 {
		table := t.GetTable(path[0], path[1], path[2])
		if table == nil || !table.expanded || len(table.columns) == 0 {
			return path
		}
		// Go to last column
		return []int{path[0], path[1], path[2], len(table.columns) - 1}
	}

	// Column level - no children
	return path
}

// SiblingCount returns the number of siblings at a specific path level
func (t *TreeState) SiblingCount(level int) int {
	switch level {
	case 0: // DatabaseLevel
		return len(t.databases)
	case 1: // SchemaLevel
		if len(t.cursor.path) > 0 {
			db := t.GetDatabase(t.cursor.path[0])
			if db != nil {
				return len(db.schemas)
			}
		}
	case 2: // TableLevel
		if len(t.cursor.path) > 1 {
			schema := t.GetSchema(t.cursor.path[0], t.cursor.path[1])
			if schema != nil {
				return len(schema.tables)
			}
		}
	case 3: // TableColumnLevel
		if len(t.cursor.path) > 2 {
			table := t.GetTable(t.cursor.path[0], t.cursor.path[1], t.cursor.path[2])
			if table != nil {
				return len(table.columns)
			}
		}
	}
	return 0
}

func (t *TreeState) Collapse() {
	switch t.cursor.Level() {
	case DatabaseLevel:
		db := t.CurrentDatabase()
		if db != nil {
			db.expanded = false
		}

	case SchemaLevel:
		schema := t.CurrentSchema()
		if schema != nil && schema.expanded {
			schema.expanded = false
		} else {
			// Schema not expanded, collapse parent database and move up
			db := t.CurrentDatabase()
			if db != nil {
				db.expanded = false
			}
			t.cursor.path = []int{t.cursor.DbIndex()}
		}

	case TableLevel:
		table := t.CurrentTable()
		if table != nil && table.expanded {
			table.expanded = false
		} else {
			// Table not expanded, collapse parent schema and move up
			schema := t.CurrentSchema()
			if schema != nil {
				schema.expanded = false
			}
			t.cursor.path = []int{t.cursor.DbIndex(), t.cursor.SchemaIndex()}
		}

	case TableColumnLevel:
		// Move to parent table and collapse it
		table := t.CurrentTable()
		if table != nil {
			table.expanded = false
		}
		t.cursor.path = []int{t.cursor.DbIndex(), t.cursor.SchemaIndex(), t.cursor.TableIndex()}
	}
}

func (t *TreeState) Expand(registry *database.DBRegistry) tea.Cmd {
	switch t.cursor.Level() {
	case DatabaseLevel:
		currentDB := t.CurrentDatabase()

		if !currentDB.parsed {
			// Connect and fetch schemas
			return handleDBSelection(t.cursor.DbIndex(), registry)
		} else if currentDB != nil {
			currentDB.expanded = true
		}

	case SchemaLevel:
		schema := t.CurrentSchema()
		if schema == nil {
			return nil
		}

		if len(schema.tables) == 0 {
			// Load tables
			return handleSchemaSelection(t.cursor.DbIndex(), t.cursor.SchemaIndex(), registry)
		} else {
			schema.expanded = true
		}

	case TableLevel:
		table := t.CurrentTable()
		if table != nil && len(table.columns) > 0 {
			table.expanded = true
		}
	}
	return nil
}

func (t *TreeState) SetSchemas(schemas []*database.Schema) {
	db := t.CurrentDatabase()
	if db == nil {
		return
	}
	db.schemas = make([]*databaseSchemaNode, 0, len(schemas))
	for _, schema := range schemas {
		db.schemas = append(db.schemas, &databaseSchemaNode{
			name:   schema.Name,
			schema: schema,
		})
	}
	db.parsed = true
	db.expanded = true
}

func (t *TreeState) SetTables(tables []*database.Table) {
	schema := t.CurrentSchema()
	if schema == nil {
		return
	}
	schema.tables = make([]*schemaTableNode, 0, len(tables))
	for _, table := range tables {
		schema.tables = append(schema.tables, &schemaTableNode{
			name:  table.Name,
			table: table,
		})
	}
	schema.expanded = true
}

func (t *TreeState) SetColumns(databaseID, schemaName string, columns map[string][]*database.Column) error {
	db, _ := t.FindDatabase(databaseID)
	if db == nil {
		return fmt.Errorf("database not found")
	}

	schema, _ := t.FindSchema(db, schemaName)
	if schema == nil {
		return fmt.Errorf("schema not found")
	}

	for _, tableNode := range schema.tables {
		for _, col := range columns[tableNode.name] {
			tableNode.columns = append(tableNode.columns, &tableColumnNode{
				name:      col.Name,
				dataType:  col.DataType,
				maxLength: col.MaxLength,
			})
		}
	}
	return nil
}

func (t *TreeState) Toggle(registry *database.DBRegistry) tea.Cmd {
	switch t.cursor.Level() {
	case DatabaseLevel:
		db := t.CurrentDatabase()
		if db == nil {
			return nil
		}
		if !db.parsed {
			return handleDBSelection(t.cursor.DbIndex(), registry)
		}
		db.expanded = !db.expanded

	case SchemaLevel:
		schema := t.CurrentSchema()
		if schema == nil {
			return nil
		}
		if len(schema.tables) == 0 {
			return handleSchemaSelection(t.cursor.DbIndex(), t.cursor.SchemaIndex(), registry)
		}
		schema.expanded = !schema.expanded

	case TableLevel:
		if table := t.CurrentTable(); table != nil {
			table.expanded = !table.expanded
		}
	}
	return nil
}
