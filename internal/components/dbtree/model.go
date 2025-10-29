package dbtree

import "github.com/SavingFrame/dbettier/internal/database"

// import tea "github.com/charmbracelet/bubbletea"

type databaseSchemaNode struct {
	name string
}

type databaseNode struct {
	name    string
	host    string
	schemas []*databaseSchemaNode
}

type DBTreeModel struct {
	focusIndex int
	databases  []*databaseNode
}

func DBTreeScreen() DBTreeModel {
	var dbNodes []*databaseNode
	for _, db := range database.Connections {
		dbNodes = append(dbNodes, &databaseNode{
			name: db.Database,
			host: db.Host,
		})
	}
	return DBTreeModel{
		focusIndex: 0,
		databases:  dbNodes,
	}
}

func (m DBTreeModel) totalFocusableItems() int {
	return len(database.Connections)
}
