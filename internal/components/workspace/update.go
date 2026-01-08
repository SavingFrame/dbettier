package workspace

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/components/logpanel"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
	zone "github.com/lrstanley/bubblezone/v2"
)

// Update handles messages for the workspace
func (w Workspace) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return w, nil
		}

		// Check if click is on a tab
		for i := range w.tabs {
			zoneID := fmt.Sprintf("tab-%d", i)
			if zone.Get(zoneID).InBounds(msg) {
				// Check if it's a close button click
				// The close button is at the right edge of the tab
				zoneInfo := zone.Get(zoneID)
				relativeX := msg.X - zoneInfo.StartX

				if w.IsCloseButtonClick(i, relativeX) {
					w.CloseTab(i)
				} else {
					w.SetActiveIndex(i)
				}
				return w, nil
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.NextTab):
			w.NextTab()
			return w, nil
		case key.Matches(msg, DefaultKeyMap.PrevTab):
			w.PrevTab()
			return w, nil
		case key.Matches(msg, DefaultKeyMap.CloseTab):
			w.CloseActiveTab()
			return w, nil
		}
	case sharedcomponents.ExecuteSQLTextMsg:
		t := w.ActiveTab()
		t.DatabaseID = msg.DatabaseID
		q := sharedcomponents.NewBasicSQLQuery(msg.Query)
		return w, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			executeSQLQuery(w.registry, q, msg.DatabaseID),
		)

	case sharedcomponents.ReapplyTableQueryMsg:
		return w, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			executeSQLQuery(w.registry, msg.Query, w.ActiveTab().DatabaseID),
		)

	case sharedcomponents.OpenTableMsg:
		w.AddTableTab(msg.Table.Name, msg.DatabaseID)
		t := w.ActiveTab()
		t.DatabaseID = msg.DatabaseID
		return w, tea.Batch(
			func() tea.Msg { return sharedcomponents.TableLoadingMsg{} },
			openTableHandler(w.registry, msg.Table, msg.DatabaseID),
		)
	}

	return w, tea.Batch(cmds...)
}

// HandleMouseClick processes mouse clicks for the workspace tab bar
// Returns true if the click was handled, along with any commands
func (w *Workspace) HandleMouseClick(msg tea.MouseReleaseMsg) (bool, tea.Cmd) {
	if msg.Button != tea.MouseLeft {
		return false, nil
	}

	// Check each tab zone
	for i := range w.tabs {
		zoneID := fmt.Sprintf("tab-%d", i)
		if zone.Get(zoneID).InBounds(msg) {
			zoneInfo := zone.Get(zoneID)
			relativeX := msg.X - zoneInfo.StartX

			if w.IsCloseButtonClick(i, relativeX) {
				w.CloseTab(i)
			} else {
				w.SetActiveIndex(i)
			}
			return true, nil
		}
	}

	return false, nil
}

// HandleKeys processes keyboard input for tab navigation
// Returns true if the key was handled
func (w *Workspace) HandleKeys(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, DefaultKeyMap.NextTab):
		w.NextTab()
		return true
	case key.Matches(msg, DefaultKeyMap.PrevTab):
		w.PrevTab()
		return true
	case key.Matches(msg, DefaultKeyMap.CloseTab):
		w.CloseActiveTab()
		return true
	}
	return false
}

// RouteToActiveTab routes a message to the active tab's components
// Returns commands from both tableview and sqlcommandbar
func (w *Workspace) RouteToActiveTab(msg tea.Msg, targets sharedcomponents.ComponentTarget) []tea.Cmd {
	var cmds []tea.Cmd

	tab := w.ActiveTab()
	if tab == nil {
		return cmds
	}

	if targets&sharedcomponents.TargetTableView != 0 {
		cmd := w.UpdateActiveTableView(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if targets&sharedcomponents.TargetSQLCommandBar != 0 {
		cmd := w.UpdateActiveSQLCommandBar(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

// FindTabByZone returns the tab index from a zone ID like "tab-0", "tab-1", etc.
func FindTabByZone(zoneID string) (int, bool) {
	if !strings.HasPrefix(zoneID, "tab-") {
		return -1, false
	}

	var idx int
	_, err := fmt.Sscanf(zoneID, "tab-%d", &idx)
	if err != nil {
		return -1, false
	}

	return idx, true
}

// TODO: Refactor this all this below:

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
