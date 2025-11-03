package tableview

import (
	"fmt"
	"log"

	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case sharedcomponents.SetSQLTextMsg:
		log.Println("Get `SETSQLTEXT` message in TableViewModel")
		// m.textarea.SetValue(msg.Command)
	case sharedcomponents.SQLResultMsg:
		m.handleSQLResultMsg(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	case tea.WindowSizeMsg:
		// Size will be handled by root screen
		return m, nil
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TableViewModel) handleSQLResultMsg(msg sharedcomponents.SQLResultMsg) {
	var columns []table.Column
	var rows []table.Row

	colWidth := (m.width - 4) / len(msg.Columns) // -4 for borders and padding
	for _, colName := range msg.Columns {
		columns = append(columns, table.Column{
			Title: colName,
			Width: colWidth,
		})
	}
	m.table.SetRows(nil)
	m.table.SetColumns(columns)
	for _, rowData := range msg.Rows {
		var rowCells []string
		for _, cell := range rowData {
			rowCells = append(rowCells, fmt.Sprintf("%v", cell))
		}
		rows = append(rows, table.Row(rowCells))
	}
	m.table.SetRows(rows)
}
