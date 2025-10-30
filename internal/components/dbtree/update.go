package dbtree

import (
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

func (m DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case handleDBSelectionResult:
		m.databases[m.cursor.dbIndex].schemas = make([]*databaseSchemaNode, 0)
		for _, schema := range msg.schemas {
			m.databases[m.cursor.dbIndex].schemas = append(m.databases[m.cursor.dbIndex].schemas, &databaseSchemaNode{
				name: schema.Name,
			})
		}
		m.databases[m.cursor.dbIndex].expanded = true
		return m, msg.notification
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			m, cmd = m.moveCursorUp()
		case "j", "down":
			m, cmd = m.moveCursorDown()
		case "enter":
			m, cmd = m.handleEnter()
		case "left", "h":
			m.collapseNode()
		case "right", "l":
			m, cmd = m.expandNode()
		}
	}
	return m, cmd
}

func (m DBTreeModel) moveCursorUp() (DBTreeModel, tea.Cmd) {
	if m.cursor.isAtDatabaseLevel() {
		// Move to previous database or its last schema if expanded
		if m.cursor.dbIndex > 0 {
			m.cursor.dbIndex--
			prevDB := m.databases[m.cursor.dbIndex]
			if prevDB.expanded && len(prevDB.schemas) > 0 {
				// Move to last schema of previous database
				m.cursor.schemaIndex = len(prevDB.schemas) - 1
			}
		} else {
			// Wrap to last item
			m.cursor.dbIndex = len(m.databases) - 1
			lastDB := m.databases[m.cursor.dbIndex]
			if lastDB.expanded && len(lastDB.schemas) > 0 {
				m.cursor.schemaIndex = len(lastDB.schemas) - 1
			}
		}
	} else {
		// Currently on a schema
		if m.cursor.schemaIndex > 0 {
			// Move to previous schema
			m.cursor.schemaIndex--
		} else {
			// Move to parent database
			m.cursor.schemaIndex = -1
		}
	}
	return m, nil
}

func (m DBTreeModel) moveCursorDown() (DBTreeModel, tea.Cmd) {
	currentDB := m.databases[m.cursor.dbIndex]

	if m.cursor.isAtDatabaseLevel() {
		// Currently on a database
		if currentDB.expanded && len(currentDB.schemas) > 0 {
			// Move to first schema
			m.cursor.schemaIndex = 0
		} else {
			// Move to next database
			if m.cursor.dbIndex < len(m.databases)-1 {
				m.cursor.dbIndex++
			} else {
				// Wrap to first database
				m.cursor.dbIndex = 0
			}
		}
	} else {
		// Currently on a schema
		if m.cursor.schemaIndex < len(currentDB.schemas)-1 {
			// Move to next schema
			m.cursor.schemaIndex++
		} else {
			// Move to next database
			m.cursor.schemaIndex = -1
			if m.cursor.dbIndex < len(m.databases)-1 {
				m.cursor.dbIndex++
			} else {
				// Wrap to first database
				m.cursor.dbIndex = 0
			}
		}
	}
	return m, nil
}

func (m DBTreeModel) handleEnter() (DBTreeModel, tea.Cmd) {
	if m.cursor.isAtDatabaseLevel() {
		// Toggle expand/collapse or connect if not connected
		currentDB := m.databases[m.cursor.dbIndex]
		dbConn := m.registry.GetAll()[m.cursor.dbIndex]

		if !dbConn.Connected {
			// Connect and fetch schemas
			return m, handleDBSelection(m.cursor.dbIndex, m.registry)
		} else {
			// Toggle expand
			currentDB.expanded = !currentDB.expanded
		}
	}
	return m, nil
}

func (m DBTreeModel) collapseNode() {
	if m.cursor.isAtDatabaseLevel() {
		// Collapse current database
		m.databases[m.cursor.dbIndex].expanded = false
	} else {
		// Move cursor to parent database and collapse
		m.cursor.schemaIndex = -1
		m.databases[m.cursor.dbIndex].expanded = false
	}
}

func (m DBTreeModel) expandNode() (DBTreeModel, tea.Cmd) {
	if m.cursor.isAtDatabaseLevel() {
		currentDB := m.databases[m.cursor.dbIndex]
		dbConn := m.registry.GetAll()[m.cursor.dbIndex]

		if !dbConn.Connected {
			// Connect and fetch schemas
			return m, handleDBSelection(m.cursor.dbIndex, m.registry)
		} else {
			// Expand
			currentDB.expanded = true
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
				// return m, nil, err
				return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
			}
		}
		schemas, err := db.ParseSchemas()
		if err != nil {
			return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
		}
		return handleDBSelectionResult{notification: notifications.ShowInfo("Successfully connected to database."), schemas: schemas}
	}
}
