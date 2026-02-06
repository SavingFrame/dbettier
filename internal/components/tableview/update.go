package tableview

import (
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/query"
	"github.com/SavingFrame/dbettier/pkgs/table"
	zone "github.com/lrstanley/bubblezone/v2"
)

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// always update spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// update status bar text input
	switch m.statusBar.focus {
	case StatusBarFocusFilter:
		cmd, earlyExit := m.handleFilterInput(msg)
		cmds = append(cmds, cmd)
		if earlyExit {
			return m, tea.Batch(cmds...)
		}
	case StatusBarFocusOrdering:
		cmd, earlyExit := m.handleOrderingInput(msg)
		cmds = append(cmds, cmd)
		if earlyExit {
			return m, tea.Batch(cmds...)
		}
	}

	// always update upstream table model
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case messages.TableLoadingMsg:
		m.isLoading = true
		return m, tea.Batch(cmds...)
	case query.SQLResultMsg:
		log.Printf("Received SQLResultMsg for TableViewModel: %+v", msg)
		m.isLoading = false
		result := m.data.SetFromSQLResult(msg)
		columns, rows := m.data.BuildTableData(result)
		m.table.SetRows(nil)
		m.table.SetColumns(columns)
		log.Println("Setting table rows")
		m.table.SetRows(rows)
		log.Println("TableViewModel update complete after SQLResultMsg")
		log.Printf("Table has %d columns and %d rows", len(m.table.Columns()), len(m.table.Rows()))
	case query.UpdateTableMsg:
		m.isLoading = false
		m.data.SetQuery(msg.Query)
		columns, rows := m.data.BuildTableData(msg.Query.GetSQLResult())
		m.table.SetRows(nil)
		m.table.SetColumns(columns)
		m.table.SetRows(rows)
	case table.SortChangeMsg:
		cmds = append(cmds, m.handleSortChange(msg))
	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return m, nil
		}
		if zone.Get("refresh").InBounds(msg) {
			cmd = m.data.RefreshQuery()
			cmds = append(cmds, cmd)
		} else if zone.Get("filterInput").InBounds(msg) {
			m.statusBar.SetFocus(StatusBarFocusFilter)
		} else if zone.Get("orderingInput").InBounds(msg) {
			m.statusBar.SetFocus(StatusBarFocusOrdering)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.statusBar.Pagination().Clear()
		case key.Matches(msg, DefaultKeyMap.Quit):
			if m.statusBar.Pagination().HasPendingConfirm() {
				m.statusBar.Pagination().Clear()
			} else {
				cmds = append(cmds, tea.Quit)
			}
		case key.Matches(msg, DefaultKeyMap.NextPage):
			cmd = m.handleNextPage()
			cmds = append(cmds, cmd)
		case key.Matches(msg, DefaultKeyMap.PreviousPage):
			cmd = m.handlePrevPage()
			cmds = append(cmds, cmd)
		default:
			m.statusBar.Pagination().Clear()
		}
	}

	// Sync status bar state before returning
	m.syncStatusBar()

	return m, tea.Batch(cmds...)
}

func (m *TableViewModel) syncStatusBar() {
	focusedRow, focusedCol := m.table.FocusedPosition()
	totalRows := len(m.table.Rows())
	totalCols := len(m.table.Columns())

	m.statusBar.SyncState(
		focusedRow,
		totalRows,
		m.data.PageOffset(),
		m.data.CanFetchTotal(),
		focusedCol,
		totalCols,
		m.data.IsTableQuery(),
		m.data.GetSortOrders(),
	)
}

func (m *TableViewModel) handleSortChange(msg table.SortChangeMsg) tea.Cmd {
	if m.data.databaseID == "" {
		log.Println("Cannot sort: no base query or database ID")
		return nil
	}

	orderByClauses := m.data.HandleSortChange(m.table.Columns(), msg.SortOrders)
	switch tq := m.data.Query().(type) {
	case *query.TableQuery:
		cmd := tq.HandleSortChange(orderByClauses)
		return cmd
	}
	return nil
}

func (m *TableViewModel) handleNextPage() tea.Cmd {
	m.table.ScrollToBottom()

	if !m.table.IsLatestRowFocused() || !m.data.HasNextPage() {
		m.statusBar.Pagination().Clear()
		return nil
	}

	if m.statusBar.Pagination().ConfirmNextPage() {
		return m.data.Query().NextPage()
	}

	m.statusBar.Pagination().RequestNextPage()
	return nil
}

func (m *TableViewModel) handlePrevPage() tea.Cmd {
	m.table.ScrollToTop()

	if !m.table.IsFirstRowFocused() || !m.data.HasPreviousPage() {
		m.statusBar.Pagination().Clear()
		return nil
	}

	// Check if this is confirmation press
	if m.statusBar.Pagination().ConfirmPrevPage() {
		return m.data.Query().PreviousPage()
	}

	// First press - request confirmation
	m.statusBar.Pagination().RequestPrevPage()
	return nil
}

// handleFilterInput processes input when the filter input is focused
// returns a command and a boolean indicating if we need to execute early exit in the update method
func (m *TableViewModel) handleFilterInput(msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.statusBar.filterInput, cmd = m.statusBar.FilterInput().Update(msg)
	cmds = append(cmds, cmd)

	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, DefaultKeyMap.Enter) {
			v := m.statusBar.FilterValue()
			if tq, ok := m.data.Query().(*query.TableQuery); ok {
				tableUpdateCmd := tq.SetWhereClause(v)
				cmds = append(cmds, tableUpdateCmd)
				return tea.Batch(cmds...), false
			}
		} else if key.Matches(msg, DefaultKeyMap.Escape) {
			m.statusBar.SetFocus(StatusBarFocusNone)
			return tea.Batch(cmds...), true
		} else {
			return tea.Batch(cmds...), true
		}
	}
	return tea.Batch(cmds...), false
}

func (m *TableViewModel) handleOrderingInput(msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.statusBar.orderingInput, cmd = m.statusBar.OrderingInput().Update(msg)
	cmds = append(cmds, cmd)

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return tea.Batch(cmds...), false
	}

	switch {
	case key.Matches(keyMsg, DefaultKeyMap.Enter):
		cmd, handled := m.applyOrdering()
		if handled {
			cmds = append(cmds, cmd)
		}
		return tea.Batch(cmds...), false
	case key.Matches(keyMsg, DefaultKeyMap.Escape):
		m.statusBar.SetFocus(StatusBarFocusNone)
		return tea.Batch(cmds...), true
	default:
		return tea.Batch(cmds...), true
	}
}

// applyOrdering parses the ordering input and applies it to the table query.
// Returns the command to execute and whether the ordering was successfully applied.
func (m *TableViewModel) applyOrdering() (tea.Cmd, bool) {
	tq, ok := m.data.Query().(*query.TableQuery)
	if !ok {
		return nil, false
	}

	orderByClauses, err := query.ParseOrderByClauses(m.statusBar.OrderingValue())
	if err != nil {
		// TODO: show error notification to user
		return nil, false
	}

	m.table.SetSortVisually(m.orderClausesToOrderCols(orderByClauses))
	return tq.HandleSortChange(orderByClauses), true
}

// orderClausesToOrderCols converts OrderByClauses to table.OrderCol slice
// by resolving column names to their indices.
func (m *TableViewModel) orderClausesToOrderCols(clauses query.OrderByClauses) []table.OrderCol {
	sortOrders := make([]table.OrderCol, 0, len(clauses))
	for _, ob := range clauses {
		sortOrders = append(sortOrders, table.OrderCol{
			Direction:   table.NewSortDirection(ob.Direction),
			ColumnIndex: m.table.ColumnIndexByName(ob.ColumnName),
		})
	}
	return sortOrders
}
