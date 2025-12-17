package tableview

import (
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/pkgs/table"
)

type TableViewModel struct {
	viewport  Viewport
	statusBar StatusBar
	data      DataState
	table     table.Model
}

func TableViewScreen() TableViewModel {
	t := table.New(
		table.WithColumns(defaultColumns()),
		table.WithRows(defaultRows()),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return TableViewModel{
		viewport:  Viewport{},
		statusBar: NewStatusBar(),
		data:      DataState{},
		table:     t,
	}
}

func (m TableViewModel) Init() tea.Cmd {
	return nil
}

func (m *TableViewModel) SetSize(width, height int) {
	m.viewport.SetSize(width, height)

	// Update table width for horizontal scrolling
	m.table.SetWidth(width)

	// Update status bar width for right-aligned content
	m.statusBar.SetWidth(width)

	// Update table height (leave some room for borders/padding and scroll indicators)
	if height > 4 {
		m.table.SetHeight(height - 4)
	}
}

func (m *TableViewModel) GetSize() (int, int) {
	return m.viewport.Width(), m.viewport.Height()
}

func (m *TableViewModel) Table() *table.Model {
	return &m.table
}
