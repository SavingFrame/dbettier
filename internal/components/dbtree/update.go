package dbtree

import (
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
)

func (m DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetSize(msg.Width, msg.Height)
		return m, nil
	case handleDBSelectionResult:
		m.tree.SetSchemas(msg.schemas)
		m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		return m, msg.notification
	case handleSchemaSelectionResult:
		m.tree.SetTables(msg.tables)
		m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		return m, msg.cmd
	case loadTablesColumnsResult:
		err := m.tree.SetColumns(msg.databaseID, msg.schemaName, msg.columns)
		if err != nil {
			return m, notifications.ShowError(err.Error())
		}
		return m, msg.notification

	case tea.KeyMsg:
		// Handle search mode input
		if m.search.IsActive() {
			m.search.HandleInput(msg, &m.tree, m.tree.cursor)
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
			return m, nil
		}
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			m.tree.cursor.MoveUp(&m.tree)
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.Down):
			m.tree.cursor.MoveDown(&m.tree)
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))

		case key.Matches(msg, DefaultKeyMap.Left):
			m.tree.Collapse()
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.Right):
			cmd = m.tree.Expand(m.registry)
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.Space):
			cmd = m.tree.Toggle(m.registry)
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.Enter):
			if m.tree.cursor.Level() != TableLevel {
				cmd = m.tree.Toggle(m.registry)
				m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
			} else {
				cmd = handleOpenTable(m.tree.CurrentDatabase(), m.tree.CurrentTable())
			}
		case key.Matches(msg, DefaultKeyMap.ScrollDown):
			for i := 0; i < m.viewport.Height()/2; i++ {
				m.tree.cursor.MoveDown(&m.tree)
			}
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.ScrollUp):
			for i := 0; i < m.viewport.Height()/2; i++ {
				m.tree.cursor.MoveUp(&m.tree)
			}
			m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
		case key.Matches(msg, DefaultKeyMap.Search):
			m.search.Enable()
		case key.Matches(msg, DefaultKeyMap.SearchNextMatch):
			if len(m.search.Matches()) > 0 {
				m.search.NextMatch(m.tree.cursor)
				m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
			}
		case key.Matches(msg, DefaultKeyMap.SearchPrevMatch):
			if len(m.search.Matches()) > 0 {
				m.search.PrevMatch(m.tree.cursor)
				m.viewport.AdjustScrollToCursor(m.tree.cursor.VisualLine(&m.tree))
			}

		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		}
	}
	return m, cmd
}

func handleDBSelection(i int, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		db := registry.GetAll()[i]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				log.Printf("Error connecting to database: %v", err)
				return notifications.ShowError(err.Error())
			}
		}
		schemas, err := db.ParseSchemas()
		if err != nil {
			return notifications.ShowError(err.Error())
		}
		return handleDBSelectionResult{schemas: schemas, notification: notifications.ShowInfo("Successfully connected to database.")}
	}
}

func handleOpenTable(db *databaseNode, t *schemaTableNode) tea.Cmd {
	return func() tea.Msg {
		t := t.table
		return sharedcomponents.OpenTableMsg{
			Table:      t,
			DatabaseID: db.id,
		}
	}
}

func handleSchemaSelection(dbIndex, schemaIndex int, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Handling schema selection for dbIndex=%d, schemaIndex=%d", dbIndex, schemaIndex)
		db := registry.GetAll()[dbIndex]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				log.Printf("Error connecting to database: %v", err)
				return notifications.ShowError(err.Error())
			}
		}
		schema := db.Schemas[schemaIndex]
		log.Printf("Loading tables for schema: %s", schema.Name)
		tables, err := schema.LoadTables()
		if err != nil {
			log.Printf("Error loading tables for schema %s: %v", schema.Name, err)
			return notifications.ShowError(err.Error())
		}
		return handleSchemaSelectionResult{
			tables: tables,
			cmd:    tea.Batch(notifications.ShowInfo("Successfully connected to database."), loadTablesColumnsCmd(schema)),
		}
	}
}

func loadTablesColumnsCmd(schema *database.Schema) tea.Cmd {
	return func() tea.Msg {
		tables, err := schema.LoadColumns()
		if err != nil {
			log.Printf("Error loading columns for schema %s: %v", schema.Name, err)
			return notifications.ShowError(err.Error())
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
