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
		isFocused := m.cursor.dbIndex() == dbIdx && m.cursor.isAtDatabaseLevel()
		if isFocused {
			dbText = focusedStyle.Render(dbText)
		} else {
			dbText = itemStyle.Render(dbText)
		}
		t.Child(dbText)

		if db.expanded && len(db.schemas) > 0 {
			// Render schemas
			schemaTree := tree.New()
			for schemaIdx, schema := range db.schemas {
				expandIndicator := ""
				if len(schema.tables) > 0 {
					if schema.expanded {
						expandIndicator = "▼ "
					} else {
						expandIndicator = "▶ "
					}
				}

				schemaText := fmt.Sprintf("%s  %s", expandIndicator, schema.name)
				isFocused := m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx
				if isFocused {
					schemaText = focusedStyle.Render(schemaText)
				} else {
					schemaText = itemStyle.Render(schemaText)
				}
				schemaTree.Child(schemaText)
				if schema.expanded && len(schema.tables) > 0 {
					// Render tables
					tableTree := tree.New()
					for tableIdx, table := range schema.tables {
						tableText := fmt.Sprintf(" %s", table.name)
						isFocused := m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.tableIndex() == tableIdx
						if isFocused {
							tableText = focusedStyle.Render(tableText)
						} else {
							tableText = itemStyle.Render(tableText)
						}
						tableTree.Child(tableText)
						if table.expanded && len(table.columns) > 0 {
							columnTree := tree.New()
							for colIdx, column := range table.columns {
								colText := fmt.Sprintf(" %s (%s)", column.name, column.dataType)
								isFocused := m.cursor.dbIndex() == dbIdx && m.cursor.schemaIndex() == schemaIdx && m.cursor.tableIndex() == tableIdx && m.cursor.tableColumnIndex() == colIdx
								if isFocused {
									colText = focusedStyle.Render(colText)
								} else {
									colText = itemStyle.Render(colText)
								}
								columnTree.Child(colText)
							}
							tableTree.Child(columnTree)
						}
					}
					schemaTree.Child(tableTree)
				}
			}
			t.Child(schemaTree)
		}
	}

	fullContent := t.String()
	// TODO: This is probably not the best way to cut off content.
	// I think we shouldn't calculate everything and then cut lines,
	// but rather calculate only what fits in the window.
	lines := strings.Split(fullContent, "\n")
	if m.scrollOffset >= len(lines) {
		m.scrollOffset = max(0, len(lines)-m.windowHeight)
	}
	if len(lines) > m.windowHeight {
		end := min(m.scrollOffset+m.windowHeight, len(lines))
		lines = lines[m.scrollOffset:end]
	}

	return strings.Join(lines, "\n")
}
