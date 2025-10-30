package dbtree

import (
	"fmt"
	"log"
	"strings"

	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func (m DBTreeModel) Init() tea.Cmd {
	return nil
}

var (
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func (m DBTreeModel) View() string {
	var b strings.Builder
	t := tree.
		Root("Databases:").
		// Child(
		// 	"Glossier",
		// 	"Fenty Beauty",
		// 	// tree.New().Child(
		// 	// 	"Gloss Bomb Universal Lip Luminizer",
		// 	// 	"Hot Cheeks Velour Blushlighter",
		// 	// ),
		// 	"Nyx",
		// 	"Mac",
		// 	"Milk",
		// ).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyleFunc(func(c tree.Children, i int) lipgloss.Style {
			if m.focusIndex == i {
				return focusedStyle
			}
			return itemStyle
		})
	for i, db := range m.databases {
		dbConn := database.Connections[i]
		var mark string
		if dbConn.Connected {
			mark = "✔"
		} else {
			mark = "✘"
		}
		t.Child(fmt.Sprintf("%s %s@%s", mark, db.name, db.host))
		// TODO: I Think we can create FlatTree in the model and use it for focusIndex
		if len(db.schemas) > 0 {
			schemaTree := tree.New().ItemStyle(itemStyle).
				ItemStyleFunc(func(c tree.Children, i int) lipgloss.Style {
					node := c.At(i)
					log.Printf("Node at index %d: %+v", i, node.Children())
					if m.focusIndex == i {
						return focusedStyle
					}
					return itemStyle
				})

			for _, schema := range db.schemas {
				schemaTree.Child(schema.name)
			}
			t.Child(schemaTree)
		}
	}

	b.WriteString(t.String())
	b.WriteString("\n\n")
	return b.String()
}
