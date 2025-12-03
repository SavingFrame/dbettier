package sqlcommandbar

import (
	"context"
	"log"

	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "charm.land/bubbletea/v2"
)

type errMsg error

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case sharedcomponents.SetSQLTextMsg:
		m.textarea.SetValue(msg.Query.Compile())
		m.databaseID = msg.DatabaseID
		return m, executeSQLCommand(m.registry, msg.Query, msg.DatabaseID)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case "ctrl+c":
			return m, tea.Quit
			// Don't auto-focus - let root handle focus
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Focus sets focus to the textarea
func (m *SQLCommandBarModel) Focus() tea.Cmd {
	return m.textarea.Focus()
}

// Blur removes focus from the textarea
func (m *SQLCommandBarModel) Blur() {
	m.textarea.Blur()
}

func executeSQLCommand(r *database.DBRegistry, q sharedcomponents.SQLQuery, databaseID string) tea.Cmd {
	log.Println("Executing SQL command on database ID:", databaseID)
	return func() tea.Msg {
		db := r.GetByID(databaseID)
		if db == nil {
			return notifications.ShowError("Database with ID " + databaseID + " not found")
		}
		conn := db.Connection
		if conn == nil {
			// TODO: TMP
			db.Connect()
			conn = db.Connection
			// return notifications.ShowError("Database connection is nil")
		}

		log.Println("Fetching rows...")
		rows, err := conn.Query(context.Background(), q.Compile())
		log.Println("After query")
		if err != nil {
			log.Printf("Failed to execute query %s", err.Error())
			return notifications.ShowError("Failed to execute query: " + err.Error())
		}
		defer rows.Close()
		fieldDescriptions := rows.FieldDescriptions()
		columnNames := make([]string, len(fieldDescriptions))
		for i, fd := range fieldDescriptions {
			columnNames[i] = string(fd.Name)
		}
		var results [][]any
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return notifications.ShowError("Failed to read row: " + err.Error())
			}
			results = append(results, values)
		}
		log.Printf("SQL command executed, retrieved %d rows\n", len(results))

		if rows.Err() != nil {
			return notifications.ShowError("Row iteration error: " + rows.Err().Error())
		}
		return sharedcomponents.SQLResultMsg{
			Columns:    columnNames,
			Rows:       results,
			Query:      q,
			DatabaseID: databaseID,
		}
	}
}
