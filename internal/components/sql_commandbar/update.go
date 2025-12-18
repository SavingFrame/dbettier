package sqlcommandbar

import (
	"context"
	"fmt"
	"log"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/components/logpanel"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
)

// KeyMap defines keybindings for the SQL command bar component
type KeyMap struct {
	Execute key.Binding
	Quit    key.Binding
}

// DefaultKeyMap returns the default keybindings for the SQL command bar
var DefaultKeyMap = KeyMap{
	Execute: key.NewBinding(
		key.WithKeys("ctrl+enter"),
		key.WithHelp("ctrl+enter", "execute SQL"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Execute, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Execute, k.Quit},
	}
}

type errMsg error

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case sharedcomponents.ExecuteSQLTextMsg:
		m.textarea.SetValue(msg.Query)
		m.databaseID = msg.DatabaseID
		q := sharedcomponents.NewBasicSQLQuery(msg.Query)
		return m, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			executeSQLQuery(m.registry, q, msg.DatabaseID),
		)
	case sharedcomponents.ReapplyTableQueryMsg:
		return m, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			executeSQLQuery(m.registry, msg.Query, m.databaseID),
		)
	case sharedcomponents.OpenTableMsg:
		m.databaseID = msg.DatabaseID
		return m, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			openTableHandler(m.registry, msg.Table, msg.DatabaseID),
		)
	case sharedcomponents.SQLResultMsg:
		m.textarea.SetValue(msg.Query.Compile())
		m.query = msg.Query
		return m, nil
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

func executeSQLQuery(r *database.DBRegistry, q sharedcomponents.QueryCompiler, databaseID string) tea.Cmd {
	return func() tea.Msg {
		db := r.GetByID(databaseID)
		if db == nil {
			return tea.BatchMsg{
				logpanel.AddLogCmd("Database with ID "+databaseID+" not found", sharedcomponents.LogError),
				notifications.ShowError("Database with ID " + databaseID + " not found"),
			}
		}
		conn := db.Connection
		if conn == nil {
			// TODO: TMP
			err := db.Connect()
			if err != nil {
				return tea.BatchMsg{
					logpanel.AddLogCmd("Failed to connect to database: "+err.Error(), sharedcomponents.LogError),
					notifications.ShowError("Failed to connect to database: " + err.Error()),
				}
			}
			conn = db.Connection
		}

		compiledQuery := q.Compile()
		log.Printf("Executing SQL query: %s\n", compiledQuery)
		startTime := time.Now()
		rows, err := conn.Query(context.Background(), compiledQuery)
		executionTime := time.Since(startTime)
		if err != nil {
			log.Printf("Failed to execute query %s", err.Error())
			return tea.BatchMsg{
				logpanel.AddLogCmd(compiledQuery, sharedcomponents.LogSQL),
				logpanel.AddLogCmd("Failed to execute query: "+err.Error(), sharedcomponents.LogError),
				notifications.ShowError("Failed to execute query: " + err.Error()),
			}
		}
		defer rows.Close()
		fieldDescriptions := rows.FieldDescriptions()
		columnNames := make([]string, len(fieldDescriptions))
		for i, fd := range fieldDescriptions {
			columnNames[i] = string(fd.Name)
		}
		var results [][]any
		fetchStart := time.Now()
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return tea.BatchMsg{
					logpanel.AddLogCmd(compiledQuery, sharedcomponents.LogSQL),
					logpanel.AddLogCmd("Failed to read row: "+err.Error(), sharedcomponents.LogError),
					notifications.ShowError("Failed to read row: " + err.Error()),
				}
			}
			results = append(results, values)
		}
		fetchingTime := time.Since(fetchStart)
		totalTime := executionTime + fetchingTime
		log.Printf("SQL command executed, retrieved %d rows\n", len(results))

		if rows.Err() != nil {
			log.Printf("Row iteration error: %s", rows.Err().Error())
			return tea.BatchMsg{
				logpanel.AddLogCmd(compiledQuery, sharedcomponents.LogSQL),
				logpanel.AddLogCmd("Row iteration error: "+rows.Err().Error(), sharedcomponents.LogError),
				notifications.ShowError("Row iteration error: " + rows.Err().Error()),
			}
		}
		return tea.BatchMsg{
			logpanel.AddLogCmd(compiledQuery, sharedcomponents.LogSQL),
			logpanel.AddLogCmd(fmt.Sprintf("Executed query in %s(execution: %s, fetching: %s), retrieved %d rows", totalTime, executionTime, fetchingTime, len(results)), sharedcomponents.LogSuccess),
			func() tea.Msg {
				return sharedcomponents.SQLResultMsg{
					Columns:    columnNames,
					Rows:       results,
					Query:      q,
					DatabaseID: databaseID,
				}
			},
		}
	}
}

func openTableHandler(r *database.DBRegistry, table *database.Table, databaseID string) tea.Cmd {
	log.Printf("Opening table %s\n", table.Name)
	return tea.Batch(
		logpanel.AddLogCmd(fmt.Sprintf("Opening table: %s", table.Name), sharedcomponents.LogInfo),
		func() tea.Msg {
			baseQuery := fmt.Sprintf("SELECT * FROM \"%s\"", table.Name)
			q := sharedcomponents.NewTableQuery(baseQuery, 500)
			return executeSQLQuery(r, q, databaseID)()
		},
	)
}
