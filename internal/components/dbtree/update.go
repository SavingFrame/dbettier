package dbtree

import (
	"fmt"
	"log"

	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Space      key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Enter      key.Binding
	Quite      key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("</h", "collapse node"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp(">/l", "expand node"),
	),
	Space: key.NewBinding(
		key.WithKeys("space", " "),
		key.WithHelp("space", "select/expand"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "page down"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "page up"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	Quite: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/ctrl+c/esc", "quit"),
	),
}

func (m DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		return m, nil
	case showNotificationCmd:
		log.Printf("Showing notification cmd: %v", msg)
		return m, msg.cmd
	case handleDBSelectionResult:
		db := m.getCurrentDatabase()
		if db != nil {
			db.schemas = make([]*databaseSchemaNode, 0)
			for _, schema := range msg.schemas {
				db.schemas = append(db.schemas, &databaseSchemaNode{
					name: schema.Name,
				})
			}
			db.parsed = true
			db.expanded = true
		}
		m = m.adjustScrollToCursor()
		return m, msg.notification
	case handleSchemaSelectionResult:
		schema := m.getCurrentSchema()
		if schema != nil {
			schema.tables = make([]*schemaTableNode, 0)
			for _, table := range msg.tables {
				schema.tables = append(schema.tables, &schemaTableNode{
					name: table.Name,
				})
			}
			schema.expanded = true
			log.Printf("Schema %s tables loaded", schema.name)
		}
		m = m.adjustScrollToCursor()
		return m, msg.cmd
	case loadTablesColumnsResult:
		database, _ := m.findDatabase(msg.databaseID)
		if database == nil {
			return m, notifications.ShowError("Database not found for loading columns.")
		}

		schema, _ := m.findSchema(database, msg.schemaName)
		if schema == nil {
			return m, notifications.ShowError("Schema not found for loading columns.")
		}

		for _, tableNode := range schema.tables {
			for _, col := range msg.columns[tableNode.name] {
				tableNode.columns = append(tableNode.columns, &tableColumnNode{
					name:      col.Name,
					dataType:  col.DataType,
					maxLength: col.MaxLength,
				})
			}
		}
		return m, msg.notification

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			m, cmd = m.moveCursorUp()
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.Down):
			m, cmd = m.moveCursorDown()
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.Space):
			m, cmd = m.handleSpace()
		case key.Matches(msg, DefaultKeyMap.Enter):
			log.Println("DBTreeModel: Enter key pressed")
			if m.cursor.level() != TableLevel {
				m, cmd = m.handleSpace()
				return m, cmd
			}
			return m, handleOpenDatabase(m.getCurrentDatabase(), m.getCurrentSchema(), m.getCurrentTable())
		case key.Matches(msg, DefaultKeyMap.Left):
			m = m.collapseNode()
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.Right):
			m, cmd = m.expandNode()
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.ScrollDown):
			scrollAmount := m.windowHeight / 2
			for range scrollAmount {
				m, _ = m.moveCursorDown()
			}
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.ScrollUp):
			scrollAmount := m.windowHeight / 2
			for range scrollAmount {
				m, _ = m.moveCursorUp()
			}
			m = m.adjustScrollToCursor()
		case key.Matches(msg, DefaultKeyMap.Quite):
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m DBTreeModel) moveCursorUp() (DBTreeModel, tea.Cmd) {
	currentIdx := m.getCurrentIndex()

	// Try to move to previous sibling
	if currentIdx > 0 {
		m.cursor.path[len(m.cursor.path)-1]--
		// Navigate to the last visible descendant of the previous sibling
		m.cursor.path = m.getLastVisibleDescendantAtPath(m.cursor.path)
	} else {
		// Move to parent level
		if len(m.cursor.path) > 1 {
			m.cursor.path = m.cursor.path[:len(m.cursor.path)-1]
		}
	}
	return m, nil
}

func (m DBTreeModel) moveCursorDown() (DBTreeModel, tea.Cmd) {
	// Try to move into children first
	if m.isExpanded() && m.hasChildren() {
		m.cursor.path = append(m.cursor.path, 0)
		return m, nil
	}

	// Try to move to next sibling at current or any parent level
	for level := len(m.cursor.path); level > 0; level-- {
		// Get the index at this level
		currentIdx := m.cursor.path[level-1]
		siblingCount := m.getSiblingCountAtLevel(level - 1)

		if currentIdx < siblingCount-1 {
			// Move to next sibling at this level
			m.cursor.path = m.cursor.path[:level]
			m.cursor.path[level-1]++
			return m, nil
		}
	}

	return m, nil
}

// getSiblingCountAtLevel returns the number of siblings at a specific path level
func (m DBTreeModel) getSiblingCountAtLevel(level int) int {
	switch level {
	case 0: // DatabaseLevel
		return len(m.databases)
	case 1: // SchemaLevel
		if len(m.cursor.path) > 0 {
			db := m.getDatabase(m.cursor.path[0])
			if db != nil {
				return len(db.schemas)
			}
		}
	case 2: // TableLevel
		if len(m.cursor.path) > 1 {
			schema := m.getSchema(m.cursor.path[0], m.cursor.path[1])
			if schema != nil {
				return len(schema.tables)
			}
		}
	case 3: // TableColumnLevel
		if len(m.cursor.path) > 2 {
			table := m.getTable(m.cursor.path[0], m.cursor.path[1], m.cursor.path[2])
			if table != nil {
				return len(table.columns)
			}
		}
	}
	return 0
}

func (m DBTreeModel) handleSpace() (DBTreeModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.cursor.level() {
	case DatabaseLevel:
		dbIdx := m.cursor.dbIndex()
		currentDB := m.getCurrentDatabase()

		if !currentDB.parsed {
			cmd = handleDBSelection(dbIdx, m.registry)
		} else if currentDB != nil {
			currentDB.expanded = !currentDB.expanded
		}

	case SchemaLevel:
		cmd = handleSchemaSelection(m.cursor.dbIndex(), m.cursor.schemaIndex(), m.registry)

	case TableLevel:
		table := m.getCurrentTable()
		if table != nil {
			table.expanded = !table.expanded
			m = m.adjustScrollToCursor()
		}
	}
	return m, cmd
}

func (m DBTreeModel) collapseNode() DBTreeModel {
	switch m.cursor.level() {
	case DatabaseLevel:
		db := m.getCurrentDatabase()
		if db != nil {
			db.expanded = false
		}

	case SchemaLevel:
		schema := m.getCurrentSchema()
		if schema != nil && schema.expanded {
			schema.expanded = false
		} else {
			// Schema not expanded, collapse parent database and move up
			db := m.getCurrentDatabase()
			if db != nil {
				db.expanded = false
			}
			m.cursor.path = []int{m.cursor.dbIndex()}
		}

	case TableLevel:
		table := m.getCurrentTable()
		if table != nil && table.expanded {
			table.expanded = false
		} else {
			// Table not expanded, collapse parent schema and move up
			schema := m.getCurrentSchema()
			if schema != nil {
				schema.expanded = false
			}
			m.cursor.path = []int{m.cursor.dbIndex(), m.cursor.schemaIndex()}
		}

	case TableColumnLevel:
		// Move to parent table and collapse it
		table := m.getCurrentTable()
		if table != nil {
			table.expanded = false
		}
		m.cursor.path = []int{m.cursor.dbIndex(), m.cursor.schemaIndex(), m.cursor.tableIndex()}
	}
	return m
}

func (m DBTreeModel) expandNode() (DBTreeModel, tea.Cmd) {
	switch m.cursor.level() {
	case DatabaseLevel:
		dbIdx := m.cursor.dbIndex()
		currentDB := m.getCurrentDatabase()

		if !currentDB.parsed {
			// Connect and fetch schemas
			return m, handleDBSelection(dbIdx, m.registry)
		} else if currentDB != nil {
			currentDB.expanded = true
		}

	case SchemaLevel:
		schema := m.getCurrentSchema()
		if schema == nil {
			return m, nil
		}

		if len(schema.tables) == 0 {
			// Load tables
			return m, handleSchemaSelection(m.cursor.dbIndex(), m.cursor.schemaIndex(), m.registry)
		} else {
			schema.expanded = true
		}

	case TableLevel:
		table := m.getCurrentTable()
		if table != nil && len(table.columns) > 0 {
			table.expanded = true
		}
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
				log.Printf("Error connecting to database: %v", err)
				return showNotificationCmd{cmd: notifications.ShowError(err.Error())}
			}
		}
		schemas, err := db.ParseSchemas()
		if err != nil {
			return showNotificationCmd{cmd: notifications.ShowError(err.Error())}
		}
		return handleDBSelectionResult{schemas: schemas, notification: notifications.ShowInfo("Successfully connected to database.")}
	}
}

func handleOpenDatabase(db *databaseNode, s *databaseSchemaNode, t *schemaTableNode) tea.Cmd {
	return func() tea.Msg {
		cmd := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 500;", s.name, t.name)
		return sharedcomponents.SetSQLTextMsg{
			Command:    cmd,
			DatabaseID: db.id,
		}
	}
}

type handleSchemaSelectionResult struct {
	tables []*database.Table
	cmd    tea.Cmd
}

func handleSchemaSelection(dbIndex, schemaIndex int, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Handling schema selection for dbIndex=%d, schemaIndex=%d", dbIndex, schemaIndex)
		db := registry.GetAll()[dbIndex]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				log.Printf("Error connecting to database: %v", err)
				return showNotificationCmd{
					cmd: notifications.ShowError(err.Error()),
				}
			}
		}
		schema := db.Schemas[schemaIndex]
		log.Printf("Loading tables for schema: %s", schema.Name)
		tables, err := schema.LoadTables()
		if err != nil {
			log.Printf("Error loading tables for schema %s: %v", schema.Name, err)
			return showNotificationCmd{
				cmd: notifications.ShowError(err.Error()),
			}
		}
		return handleSchemaSelectionResult{
			tables: tables,
			cmd:    tea.Batch(notifications.ShowInfo("Successfully connected to database."), loadTablesColumnsCmd(schema)),
		}
	}
}

type loadTablesColumnsResult struct {
	columns      map[string][]*database.Column
	schemaName   string
	databaseID   string
	notification tea.Cmd
}

func loadTablesColumnsCmd(schema *database.Schema) tea.Cmd {
	return func() tea.Msg {
		tables, err := schema.LoadColumns()
		if err != nil {
			log.Printf("Error loading columns for schema %s: %v", schema.Name, err)
			return showNotificationCmd{
				cmd: notifications.ShowError(err.Error()),
			}
		}
		tableMap := make(map[string][]*database.Column)
		for t, cols := range tables {
			tableMap[t.Name] = cols
		}
		return loadTablesColumnsResult{
			columns:    tableMap,
			databaseID: schema.Database.ID,
			schemaName: schema.Name,

			notification: notifications.ShowSuccess("Tables and columns loaded successfully."),
		}
	}
}

type showNotificationCmd struct {
	cmd tea.Cmd
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
