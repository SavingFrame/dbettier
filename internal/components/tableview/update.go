package tableview

import (
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/pkgs/table"
)

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// always update upstream table model
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case sharedcomponents.SQLResultMsg:
		result := m.data.SetFromSQLResult(msg)
		columns, rows := m.data.BuildTableData(result)
		m.table.SetRows(nil) // TODO: WHY?
		m.table.SetColumns(columns)
		m.table.SetRows(rows)
	case sharedcomponents.UpdateTableMsg:
		m.data.SetQuery(msg.Query)
		columns, rows := m.data.BuildTableData(msg.Query.GetSQLResult())
		m.table.SetRows(nil) // TODO: WHY?
		m.table.SetColumns(columns)
		m.table.SetRows(rows)
	case table.SortChangeMsg:
		cmds = append(cmds, m.handleSortChange(msg))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.pagination.Clear()
		case key.Matches(msg, DefaultKeyMap.Quit):
			if m.pagination.HasPendingConfirm() {
				m.pagination.Clear()
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
			m.pagination.Clear()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *TableViewModel) handleSortChange(msg table.SortChangeMsg) tea.Cmd {
	if m.data.databaseID == "" {
		log.Println("Cannot sort: no base query or database ID")
		return nil
	}

	orderByClauses := m.data.HandleSortChange(m.table.Columns(), msg.SortOrders)
	switch tq := m.data.Query().(type) {
	case *sharedcomponents.TableQuery:
		tq.HandleSortChange(orderByClauses)
		return func() tea.Msg {
			return sharedcomponents.ReapplyTableQueryMsg{
				Query: tq,
			}
		}
	}
	return nil
}

func (m *TableViewModel) handleNextPage() tea.Cmd {
	m.table.ScrollToBottom()

	if !m.table.IsLatestRowFocused() || !m.data.HasNextPage() {
		m.pagination.Clear()
		return nil
	}

	if m.pagination.ConfirmNextPage() {
		return m.data.Query().NextPage()
	}

	m.pagination.RequestNextPage()
	return nil
}

func (m *TableViewModel) handlePrevPage() tea.Cmd {
	m.table.ScrollToTop()

	if !m.table.IsFirstRowFocused() || !m.data.HasPreviousPage() {
		m.pagination.Clear()
		return nil
	}

	// Check if this is confirmation press
	if m.pagination.ConfirmPrevPage() {
		return m.data.Query().PreviousPage()
	}

	// First press - request confirmation
	m.pagination.RequestPrevPage()
	return nil
}
