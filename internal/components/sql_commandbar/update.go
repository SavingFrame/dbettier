package sqlcommandbar

import (
	"context"
	"log"

	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case sharedcomponents.SetSQLTextMsg:
		log.Println("Get `SETSQLTEXT` message in SQLCommandBarModel")
		m.textarea.SetValue(msg.Command)
		return m, executeSQLCommand(m.registry, msg.Command, msg.DatabaseID)
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

func executeSQLCommand(r *database.DBRegistry, command string, databaseID string) tea.Cmd {
	return func() tea.Msg {
		db := r.GetByID(databaseID)
		if db == nil {
			return notifications.ShowError("Database with ID " + databaseID + " not found")
		}
		conn := db.Connection

		rows, err := conn.Query(context.Background(), command)
		fieldDescriptions := rows.FieldDescriptions()
		columnNames := make([]string, len(fieldDescriptions))
		for i, fd := range fieldDescriptions {
			columnNames[i] = string(fd.Name)
		}
		var results [][]interface{}
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return notifications.ShowError("Failed to read row: " + err.Error())
			}
			results = append(results, values)
		}

		if rows.Err() != nil {
			return notifications.ShowError("Row iteration error: " + err.Error())
		}
		return sharedcomponents.SQLResultMsg{
			Columns: columnNames,
			Rows:    results,
		}
	}
}
