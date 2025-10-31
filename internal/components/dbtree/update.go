package dbtree

import (
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

func (m DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		return m, nil
	case handleDBSelectionResult:
		dbIdx := m.cursor.dbIndex()
		m.databases[dbIdx].schemas = make([]*databaseSchemaNode, 0)
		for _, schema := range msg.schemas {
			m.databases[dbIdx].schemas = append(m.databases[dbIdx].schemas, &databaseSchemaNode{
				name: schema.Name,
			})
		}
		m.databases[dbIdx].expanded = true
		m = m.adjustScrollToCursor()
		return m, msg.notification
	case handleSchemaSelectionResult:
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		s := m.databases[dbIdx].schemas[schemaIdx]
		s.tables = make([]*schemaTableNode, 0)
		for _, table := range msg.tables {
			s.tables = append(s.tables, &schemaTableNode{
				name: table.Name,
			})
		}
		s.expanded = true
		m = m.adjustScrollToCursor()
		return m, msg.notification
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			m, cmd = m.moveCursorUp()
			m = m.adjustScrollToCursor()
		case "j", "down":
			m, cmd = m.moveCursorDown()
			m = m.adjustScrollToCursor()
		case "enter":
			m, cmd = m.handleEnter()
		case "left", "h":
			m.collapseNode()
			m = m.adjustScrollToCursor()
		case "right", "l":
			m, cmd = m.expandNode()
			m = m.adjustScrollToCursor()
		case "ctrl+d": // Page down
			scrollAmount := m.windowHeight / 2
			for range scrollAmount {
				m, _ = m.moveCursorDown()
			}
			m = m.adjustScrollToCursor()
		case "ctrl+u": // Page up
			scrollAmount := m.windowHeight / 2
			for range scrollAmount {
				m, _ = m.moveCursorUp()
			}
			m = m.adjustScrollToCursor()
		}
	}
	return m, cmd
}

func (m DBTreeModel) moveCursorUp() (DBTreeModel, tea.Cmd) {
	switch m.cursor.level() {
	case DatabaseLevel:
		// At database level
		dbIdx := m.cursor.dbIndex()
		if dbIdx > 0 {
			// Move to previous database or its last visible descendant
			dbIdx--
			m.cursor.path = m.getLastVisibleDescendant(dbIdx)
		} else {
			// Wrap to last database and its last descendant
			dbIdx = len(m.databases) - 1
			m.cursor.path = m.getLastVisibleDescendant(dbIdx)
		}

	case SchemaLevel:
		// At schema level
		schemaIdx := m.cursor.schemaIndex()
		if schemaIdx > 0 {
			// Move to previous schema or its last visible descendant
			dbIdx := m.cursor.dbIndex()
			schemaIdx--
			m.cursor.path = []int{dbIdx, schemaIdx}
			// Check if previous schema has expanded tables
			schema := m.databases[dbIdx].schemas[schemaIdx]
			if schema.expanded && len(schema.tables) > 0 {
				m.cursor.path = append(m.cursor.path, len(schema.tables)-1)
			}
		} else {
			// Move to parent database
			m.cursor.path = []int{m.cursor.dbIndex()}
		}

	case TableLevel:
		// At table level
		tableIdx := m.cursor.tableIndex()
		if tableIdx > 0 {
			// Move to previous table
			m.cursor.path[2]--
		} else {
			// Move to parent schema
			m.cursor.path = []int{m.cursor.dbIndex(), m.cursor.schemaIndex()}
		}
	}
	return m, nil
}

func (m DBTreeModel) moveCursorDown() (DBTreeModel, tea.Cmd) {
	switch m.cursor.level() {
	case DatabaseLevel:
		// At database level
		dbIdx := m.cursor.dbIndex()
		currentDB := m.databases[dbIdx]

		if currentDB.expanded && len(currentDB.schemas) > 0 {
			// Move to first schema
			m.cursor.path = []int{dbIdx, 0}
		} else {
			// Move to next database
			if dbIdx < len(m.databases)-1 {
				m.cursor.path = []int{dbIdx + 1}
			} else {
				// Wrap to first database
				m.cursor.path = []int{0}
			}
		}

	case SchemaLevel:
		// At schema level
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		currentSchema := m.databases[dbIdx].schemas[schemaIdx]

		if currentSchema.expanded && len(currentSchema.tables) > 0 {
			// Move to first table
			m.cursor.path = []int{dbIdx, schemaIdx, 0}
		} else if schemaIdx < len(m.databases[dbIdx].schemas)-1 {
			// Move to next schema
			m.cursor.path = []int{dbIdx, schemaIdx + 1}
		} else {
			// Move to next database
			if dbIdx < len(m.databases)-1 {
				m.cursor.path = []int{dbIdx + 1}
			} else {
				// Wrap to first database
				m.cursor.path = []int{0}
			}
		}

	case TableLevel:
		// At table level
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		tableIdx := m.cursor.tableIndex()
		currentSchema := m.databases[dbIdx].schemas[schemaIdx]

		if tableIdx < len(currentSchema.tables)-1 {
			// Move to next table
			m.cursor.path[2]++
		} else if schemaIdx < len(m.databases[dbIdx].schemas)-1 {
			// Move to next schema
			m.cursor.path = []int{dbIdx, schemaIdx + 1}
		} else {
			// Move to next database
			if dbIdx < len(m.databases)-1 {
				m.cursor.path = []int{dbIdx + 1}
			} else {
				// Wrap to first database
				m.cursor.path = []int{0}
			}
		}
	}
	return m, nil
}

// getLastVisibleDescendant returns the path to the last visible descendant of a database
func (m DBTreeModel) getLastVisibleDescendant(dbIdx int) []int {
	db := m.databases[dbIdx]
	if !db.expanded || len(db.schemas) == 0 {
		return []int{dbIdx}
	}

	// Find last schema
	schemaIdx := len(db.schemas) - 1
	schema := db.schemas[schemaIdx]

	if !schema.expanded || len(schema.tables) == 0 {
		return []int{dbIdx, schemaIdx}
	}

	// Find last table
	tableIdx := len(schema.tables) - 1
	return []int{dbIdx, schemaIdx, tableIdx}
}

func (m DBTreeModel) handleEnter() (DBTreeModel, tea.Cmd) {
	switch m.cursor.level() {
	case DatabaseLevel:
		// Toggle expand/collapse or connect if not connected
		dbIdx := m.cursor.dbIndex()
		currentDB := m.databases[dbIdx]
		dbConn := m.registry.GetAll()[dbIdx]

		if !dbConn.Connected {
			// Connect and fetch schemas
			return m, handleDBSelection(dbIdx, m.registry)
		} else {
			// Toggle expand
			currentDB.expanded = !currentDB.expanded
		}

	case SchemaLevel:
		// Load tables for schema
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		return m, handleSchemaSelection(dbIdx, schemaIdx, m.registry)

	case TableLevel:
		// TODO: Handle table selection (show columns, etc.)
		// For now, do nothing
	}
	return m, nil
}

func (m DBTreeModel) collapseNode() {
	switch m.cursor.level() {
	case DatabaseLevel:
		// Collapse current database
		dbIdx := m.cursor.dbIndex()
		m.databases[dbIdx].expanded = false

	case SchemaLevel:
		// Move cursor to parent database and collapse
		dbIdx := m.cursor.dbIndex()
		m.cursor.path = []int{dbIdx}
		m.databases[dbIdx].expanded = false

	case TableLevel:
		// Move cursor to parent schema and collapse
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		m.cursor.path = []int{dbIdx, schemaIdx}
		m.databases[dbIdx].schemas[schemaIdx].expanded = false
	}
}

func (m DBTreeModel) expandNode() (DBTreeModel, tea.Cmd) {
	switch m.cursor.level() {
	case DatabaseLevel:
		dbIdx := m.cursor.dbIndex()
		currentDB := m.databases[dbIdx]
		dbConn := m.registry.GetAll()[dbIdx]

		if !dbConn.Connected {
			// Connect and fetch schemas
			return m, handleDBSelection(dbIdx, m.registry)
		} else {
			// Expand
			currentDB.expanded = true
		}

	case SchemaLevel:
		// Expand schema and load tables if not loaded
		dbIdx := m.cursor.dbIndex()
		schemaIdx := m.cursor.schemaIndex()
		schema := m.databases[dbIdx].schemas[schemaIdx]

		if len(schema.tables) == 0 {
			// Load tables
			return m, handleSchemaSelection(dbIdx, schemaIdx, m.registry)
		} else {
			schema.expanded = true
		}

	case TableLevel:
		// TODO: Expand table to show columns, indexes, etc.
		// For now, do nothing
	}
	return m, nil
}

type handleDBSelectionResult struct {
	notification tea.Cmd
	schemas      []*database.Schema
}

func handleDBSelection(i int, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		db := registry.GetAll()[i]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				// return m, nil, err
				return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
			}
		}
		schemas, err := db.ParseSchemas()
		if err != nil {
			return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
		}
		return handleDBSelectionResult{schemas: schemas, notification: notifications.ShowInfo("Successfully connected to database.")}
	}
}

type handleSchemaSelectionResult struct {
	notification tea.Cmd
	tables       []*database.Table
}

func handleSchemaSelection(dbIndex, schemaIndex int, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		db := registry.GetAll()[dbIndex]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				// return m, nil, err
				return handleSchemaSelectionResult{
					notification: notifications.ShowError(err.Error()),
				}
			}
		}
		schema := db.Schemas[schemaIndex]
		tables, err := schema.LoadTables()
		if err != nil {
			return handleSchemaSelectionResult{
				notification: notifications.ShowError(err.Error()),
			}
		}
		return handleSchemaSelectionResult{
			tables:       tables,
			notification: notifications.ShowInfo("Successfully connected to database."),
		}
	}
}

func (m DBTreeModel) adjustScrollToCursor() DBTreeModel {
	cursorLine := m.getCursorVisualLine()
	visibleHeight := m.windowHeight - 3

	// If cursor is above viewport, scroll up
	if cursorLine < m.scrollOffset {
		m.scrollOffset = cursorLine
	}

	// If cursor is below viewport, scroll down
	if cursorLine >= m.scrollOffset+visibleHeight {
		m.scrollOffset = cursorLine - visibleHeight + 1
	}

	// Ensure scrollOffset doesn't go negative
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	return m
}
