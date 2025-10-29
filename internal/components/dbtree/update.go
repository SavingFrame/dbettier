package dbtree

import (
	"fmt"

	"github.com/SavingFrame/dbettier/internal/components/notifications"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

func (m DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "k", "up", "j", "down":
			m, cmd = m.handleNavigation(msg)
		case "enter":
			var err error
			var notification tea.Cmd
			m, notification, err = m.handleDBSelection(m.focusIndex)
			if err != nil {
				cmd = notifications.ShowInfo(err.Error())
			} else {
				cmd = notification
			}
		}
	}
	return m, cmd
}

func (m DBTreeModel) handleNavigation(msg tea.KeyMsg) (DBTreeModel, tea.Cmd) {
	s := msg.String()
	if s == "k" || s == "up" {
		m.focusIndex--
	} else if s == "j" || s == "down" {
		m.focusIndex++
	}
	if m.focusIndex > m.totalFocusableItems() {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = m.totalFocusableItems()
	}
	return m, notifications.ShowInfo(fmt.Sprintf("Focused on item %d", m.focusIndex))
}

func (m DBTreeModel) handleDBSelection(i int) (DBTreeModel, tea.Cmd, error) {
	db := database.Connections[i]
	if !db.Connected {
		err := db.Connect()
		if err != nil {
			return m, nil, err
		}
	}
	schemas, err := db.ParseSchemas()
	if err != nil {
		return m, nil, err
	}
	m.databases[i].schemas = make([]*databaseSchemaNode, 0)
	for _, schema := range schemas {
		m.databases[i].schemas = append(m.databases[i].schemas, &databaseSchemaNode{
			name: schema.Name,
		})
	}
	return m, notifications.ShowInfo("Successfully connected to database."), nil
}
