package dbtree

// TreeCursor represents the current focus position in the tree using a path
// path[0] = database index
// path[1] = schema index (if at schema level or deeper)
// path[2] = table index (if at table level or deeper)
type TreeCursor struct {
	path []int
}

// Level returns the current tree Level
func (c *TreeCursor) Level() TreeLevel {
	return TreeLevel(len(c.path) - 1)
}

// AtLevel checks if cursor is at a specific level
func (c *TreeCursor) AtLevel(level TreeLevel) bool {
	return c.Level() == level
}

// DbIndex returns the database index
func (c *TreeCursor) DbIndex() int {
	if len(c.path) > 0 {
		return c.path[0]
	}
	return 0
}

// SchemaIndex returns the schema index, or -1 if not at schema level or deeper
func (c *TreeCursor) SchemaIndex() int {
	if len(c.path) > 1 {
		return c.path[1]
	}
	return -1
}

// TableIndex returns the table index, or -1 if not at table level
func (c *TreeCursor) TableIndex() int {
	if len(c.path) > 2 {
		return c.path[2]
	}
	return -1
}

func (c *TreeCursor) TableColumnIndex() int {
	if len(c.path) > 3 {
		return c.path[3]
	}
	return -1
}

func (c *TreeCursor) SetPath(path []int) {
	c.path = make([]int, len(path))
	copy(c.path, path)
}

// TODO: Remove
// isAtDatabaseLevel returns true if cursor is on a database (not a schema)
func (c *TreeCursor) isAtDatabaseLevel() bool {
	return c.AtLevel(DatabaseLevel)
}

// TODO: Refactor to avoid code duplication with rendering logic
func (c *TreeCursor) VisualLine(tree *TreeState) int {
	lineNum := 1

	for dbIdx, db := range tree.databases {
		// Current database is at lineNum
		if c.DbIndex() == dbIdx && c.isAtDatabaseLevel() {
			return lineNum
		}
		lineNum++

		// If database is expanded, count schemas
		if db.expanded && len(db.schemas) > 0 {
			for schemaIdx, schema := range db.schemas {
				if c.DbIndex() == dbIdx && c.SchemaIndex() == schemaIdx && c.AtLevel(SchemaLevel) {
					return lineNum
				}
				lineNum++

				// If schema is expanded, count tables
				if schema.expanded && len(schema.tables) > 0 {
					for tableIdx, table := range schema.tables {
						if c.DbIndex() == dbIdx && c.SchemaIndex() == schemaIdx && c.TableIndex() == tableIdx {
							return lineNum
						}
						lineNum++

						if table.expanded && len(table.columns) > 0 {
							// If table is expanded, count columns
							for columnIdx := range table.columns {
								if c.DbIndex() == dbIdx && c.SchemaIndex() == schemaIdx && c.TableIndex() == tableIdx && c.TableColumnIndex() == columnIdx {
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

func (c *TreeCursor) MoveUp(tree *TreeState) {
	currentIdx := c.CurrentIndex()

	// Try to move to previous sibling
	if currentIdx > 0 {
		c.path[len(c.path)-1]--
		// Navigate to the last visible descendant of the previous sibling
		c.path = tree.LastVisibleDescendant(c.path)
	} else {
		// Move to parent level
		if len(c.path) > 1 {
			c.path = c.path[:len(c.path)-1]
		}
	}
}

func (c *TreeCursor) MoveDown(tree *TreeState) {
	// Try to move into children first
	if tree.IsExpanded() && tree.HasChildren() {
		c.path = append(c.path, 0)
		return
	}

	// Try to move to next sibling at current or any parent level
	for level := len(c.path); level > 0; level-- {
		// Get the index at this level
		currentIdx := c.path[level-1]
		siblingCount := tree.SiblingCount(level - 1)

		if currentIdx < siblingCount-1 {
			// Move to next sibling at this level
			c.path = c.path[:level]
			c.path[level-1]++
		}
	}
}
