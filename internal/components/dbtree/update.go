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
	case handleDBSelectionResult:
		m.databases[m.focusIndex].schemas = make([]*databaseSchemaNode, 0)
		var flatNodes []*flatTreeNode
		for _, schema := range msg.schemas {
			m.databases[m.focusIndex].schemas = append(m.databases[m.focusIndex].schemas, &databaseSchemaNode{
				name: schema.Name,
			})
			flatNodes = append(flatNodes, &flatTreeNode{
				name:       schema.Name,
				typeOfNode: "schema",
			})
		}
		m.flatNodes = insertNodesAfter(m.flatNodes, m.focusIndex, flatNodes)
		return m, msg.notification
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "k", "up", "j", "down":
			m, cmd = m.handleNavigation(msg)
		case "enter":
			return m, handleDBSelection(m.focusIndex)
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

type handleDBSelectionResult struct {
	notification tea.Cmd
	schemas      []*database.Schema
}

func handleDBSelection(i int) tea.Cmd {
	return func() tea.Msg {
		db := database.Connections[i]
		if !db.Connected {
			err := db.Connect()
			if err != nil {
				// return m, nil, err
				return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
			}
		}
		schemas, err := db.ParseSchemas()
		if err != nil {
			return handleDBSelectionResult{notification: notifications.ShowError(err.Error())}
		}
		return handleDBSelectionResult{notification: notifications.ShowInfo("Successfully connected to database."), schemas: schemas}
	}
}

func insertNodesAfter(slice []*flatTreeNode, index int, nodes []*flatTreeNode) []*flatTreeNode {
	insertPos := index + 1
	result := make([]*flatTreeNode, 0, len(slice)+len(nodes))
	result = append(result, slice[:insertPos]...)
	result = append(result, nodes...)
	result = append(result, slice[insertPos:]...)
	return result
}
