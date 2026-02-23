// Package tableview provides a data table view component for displaying query results.
package tableview

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/pkgs/table"
)

type TableViewModel struct {
	viewport  Viewport
	statusBar StatusBar
	data      DataState
	table     table.Model
	spinner   spinner.Model
	isLoading bool
}

func TableViewScreen() TableViewModel {
	t := table.New(
		table.WithColumns(defaultColumns()),
		table.WithRows(defaultRows()),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle()

	return TableViewModel{
		viewport:  Viewport{},
		statusBar: NewStatusBar(),
		data:      DataState{},
		table:     t,
		spinner:   s,
		isLoading: false,
	}
}

func (m TableViewModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *TableViewModel) SetSize(width, height int) {
	m.viewport.SetSize(width, height)

	// Update table width for horizontal scrolling
	m.table.SetWidth(width)

	// Update status bar width for right-aligned content
	m.statusBar.SetWidth(width)

	// Fill the available tableview area: table content + 1 status bar row.
	// The surrounding border is already removed by the parent (borderedInner).
	tableHeight := max(1, height-1)
	m.table.SetHeight(tableHeight)
}

func (m *TableViewModel) GetSize() (int, int) {
	return m.viewport.Width(), m.viewport.Height()
}

func (m *TableViewModel) Table() *table.Model {
	return &m.table
}
