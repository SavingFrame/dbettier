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

type flatTreeNode struct {
	typeOfNode string // "database" or "schema"
	name       string
}

type DBTreeModel struct {
	focusIndex int
	databases  []*databaseNode
	flatNodes  []*flatTreeNode
}

func DBTreeScreen() DBTreeModel {
	var dbNodes []*databaseNode
	var flatTree []*flatTreeNode
	for _, db := range database.Connections {
		dbNodes = append(dbNodes, &databaseNode{
			name: db.Database,
			host: db.Host,
		})
		flatTree = append(flatTree, &flatTreeNode{
			typeOfNode: "database",
			name:       db.Database,
		})
	}
	return DBTreeModel{
		focusIndex: 0,
		databases:  dbNodes,
		flatNodes:  flatTree,
	}
}

// TODO: Use flatTree for focusIndex management
func (m DBTreeModel) totalFocusableItems() int {
	// totalFocusableSchemas := 0
	// for _, db := range m.databases {
	// 	totalFocusableSchemas += len(db.schemas)
	// }
	// return len(database.Connections) + totalFocusableSchemas
	return len(m.flatNodes)
}
