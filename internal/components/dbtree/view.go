package dbtree

import (
	"fmt"
	"strings"

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
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle)

	for dbIdx, db := range m.databases {
		dbConn := m.registry.GetAll()[dbIdx]
		var mark string
		if dbConn.Connected {
			mark = "✔"
		} else {
			mark = "✘"
		}

		// Determine expand/collapse indicator
		expandIndicator := ""
		if len(db.schemas) > 0 {
			if db.expanded {
				expandIndicator = "▼ "
			} else {
				expandIndicator = "▶ "
			}
		}

		// Render database with the correct style based on cursor
		dbText := fmt.Sprintf("%s%s %s@%s", expandIndicator, mark, db.name, db.host)
		isFocused := m.cursor.dbIndex == dbIdx && m.cursor.isAtDatabaseLevel()
		if isFocused {
			dbText = focusedStyle.Render(dbText)
		} else {
			dbText = itemStyle.Render(dbText)
		}
		t.Child(dbText)

		// Render schemas only if expanded
		if db.expanded && len(db.schemas) > 0 {
			schemaTree := tree.New()
			for schemaIdx, schema := range db.schemas {
				schemaText := schema.name
				isFocused := m.cursor.dbIndex == dbIdx && m.cursor.schemaIndex == schemaIdx
				if isFocused {
					schemaText = focusedStyle.Render(schemaText)
				} else {
					schemaText = itemStyle.Render(schemaText)
				}
				schemaTree.Child(schemaText)
			}
			t.Child(schemaTree)
		}
	}

	b.WriteString(t.String())
	b.WriteString("\n\n")
	return b.String()
}
