package tableview

import (
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
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
	if m.statusBar.focus == StatusBarFocusFilter {
		cmd, earlyExit := m.handleFilterInput(msg)
		log.Printf("Filter input handled, earlyExit=%v\n", earlyExit)
		cmds = append(cmds, cmd)
		if earlyExit {
			return m, tea.Batch(cmds...)
		}
	}

	// always update upstream table model
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case sharedcomponents.TableLoadingMsg:
		m.isLoading = true
		return m, tea.Batch(cmds...)
	case sharedcomponents.SQLResultMsg:
		m.isLoading = false
		result := m.data.SetFromSQLResult(msg)
		columns, rows := m.data.BuildTableData(result)
		m.table.SetRows(nil) // TODO: WHY?
		m.table.SetColumns(columns)
		m.table.SetRows(rows)
	case sharedcomponents.UpdateTableMsg:
		m.isLoading = false
		m.data.SetQuery(msg.Query)
		columns, rows := m.data.BuildTableData(msg.Query.GetSQLResult())
		m.table.SetRows(nil) // TODO: WHY?
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
	case *sharedcomponents.TableQuery:
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
			v := m.statusBar.FilterInput().Value()
			if tq, ok := m.data.Query().(*sharedcomponents.TableQuery); ok {
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
