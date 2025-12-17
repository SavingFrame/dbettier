package tableview

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

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
