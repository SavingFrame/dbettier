package tableview

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m TableViewModel) Update(msg tea.Msg) (TableViewModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
