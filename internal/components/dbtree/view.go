package dbtree

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func (m DBTreeModel) Init() tea.Cmd {
	return nil
}

// renderNode is a recursive helper that adds a node and its children to the tree
func (m DBTreeModel) renderNode(t *tree.Tree, nodeText string, isFocused bool, isSearchMatch bool, isActiveMatch bool, children func(*tree.Tree)) {
	var styledText string

	if isActiveMatch {
		styledText = searchMatchActiveStyle().Render(nodeText)
	} else if isSearchMatch {
		styledText = searchMatchStyle().Render(nodeText)
	} else if isFocused {
		styledText = focusedStyle().Render(nodeText)
	} else {
		styledText = itemStyle().Render(nodeText)
	}

	if children == nil {
		t.Child(styledText)
	} else {
		childTree := tree.New()
		children(childTree)
		t.Child(styledText).Child(childTree)
	}
}

// truncateText safely truncates text to fit within maxWidth, accounting for visual width
// and multi-byte characters. Returns the original text if it fits, or truncated text with "..."
func (m DBTreeModel) truncateText(text string, maxWidth int) string {
	textWidth := lipgloss.Width(text)
	if textWidth <= maxWidth {
		return text
	}

	// Convert to runes for safe slicing with multi-byte characters
	runes := []rune(text)
	result := ""
	currentWidth := 0

	for _, r := range runes {
		charStr := string(r)
		charWidth := lipgloss.Width(charStr)

		if currentWidth+charWidth+3 > maxWidth { // +3 for "..."
			break
		}

		result += charStr
		currentWidth += charWidth
	}

	return result + "..."
}

// renderDatabase adds a database node and recursively its schemas
func (m DBTreeModel) renderDatabase(t *tree.Tree, dbIdx int, db *databaseNode) {
	dbConn := m.registry.GetAll()[dbIdx]
	mark := "✘"
	if dbConn.Connected {
		mark = "✔"
	}

	expandIndicator := ""
	if len(db.schemas) > 0 {
		if db.expanded {
			expandIndicator = "▼ "
		} else {
			expandIndicator = "▶ "
		}
	}

	dbText := fmt.Sprintf("%s%s  %s@%s", expandIndicator, mark, db.name, db.host)
	dbText = m.truncateText(dbText, m.viewport.Width()-4)

	isFocused := m.tree.cursor.DbIndex() == dbIdx && m.tree.cursor.isAtDatabaseLevel()
	isSearchMatch, isActiveMatch := m.search.IsMatch([]int{dbIdx})

	// Define children renderer or nil
	var childrenFn func(*tree.Tree)
	if db.expanded && len(db.schemas) > 0 {
		childrenFn = func(childTree *tree.Tree) {
			for schemaIdx, schema := range db.schemas {
				m.renderSchema(childTree, dbIdx, schemaIdx, schema)
			}
		}
	}

	m.renderNode(t, dbText, isFocused, isSearchMatch, isActiveMatch, childrenFn)
}

// renderSchema adds a schema node and recursively its tables
func (m DBTreeModel) renderSchema(t *tree.Tree, dbIdx, schemaIdx int, schema *databaseSchemaNode) {
	expandIndicator := ""
	if len(schema.tables) > 0 {
		if schema.expanded {
			expandIndicator = "▼ "
		} else {
			expandIndicator = "▶ "
		}
	}

	schemaText := fmt.Sprintf("%s 󰑒 %s", expandIndicator, schema.name)
	schemaText = m.truncateText(schemaText, m.viewport.Width()-12)
	isFocused := m.tree.cursor.DbIndex() == dbIdx && m.tree.cursor.SchemaIndex() == schemaIdx && m.tree.cursor.AtLevel(SchemaLevel)
	isSearchMatch, isActiveMatch := m.search.IsMatch([]int{dbIdx, schemaIdx})

	var childrenFn func(*tree.Tree)
	if schema.expanded && len(schema.tables) > 0 {
		childrenFn = func(childTree *tree.Tree) {
			for tableIdx, table := range schema.tables {
				m.renderTable(childTree, dbIdx, schemaIdx, tableIdx, table)
			}
		}
	}

	m.renderNode(t, schemaText, isFocused, isSearchMatch, isActiveMatch, childrenFn)
}

// renderTable adds a table node and recursively its columns
func (m DBTreeModel) renderTable(t *tree.Tree, dbIdx, schemaIdx, tableIdx int, table *schemaTableNode) {
	tableText := fmt.Sprintf(" %s", table.name)
	tableText = m.truncateText(tableText, m.viewport.Width()-14)
	isFocused := m.tree.cursor.DbIndex() == dbIdx &&
		m.tree.cursor.SchemaIndex() == schemaIdx &&
		m.tree.cursor.TableIndex() == tableIdx &&
		m.tree.cursor.AtLevel(TableLevel)
	isSearchMatch, isActiveMatch := m.search.IsMatch([]int{dbIdx, schemaIdx, tableIdx})

	var childrenFn func(*tree.Tree)
	if table.expanded && len(table.columns) > 0 {
		childrenFn = func(childTree *tree.Tree) {
			for colIdx, column := range table.columns {
				m.renderColumn(childTree, dbIdx, schemaIdx, tableIdx, colIdx, column)
			}
		}
	}

	m.renderNode(t, tableText, isFocused, isSearchMatch, isActiveMatch, childrenFn)
}

// renderColumn adds a column node (leaf node)
func (m DBTreeModel) renderColumn(t *tree.Tree, dbIdx, schemaIdx, tableIdx, colIdx int, column *tableColumnNode) {
	colText := fmt.Sprintf("󰠵 %s (%s)", column.name, column.dataType)
	colText = m.truncateText(colText, m.viewport.Width()-18)
	isFocused := m.tree.cursor.DbIndex() == dbIdx &&
		m.tree.cursor.SchemaIndex() == schemaIdx &&
		m.tree.cursor.TableIndex() == tableIdx &&
		m.tree.cursor.TableColumnIndex() == colIdx
	isSearchMatch, isActiveMatch := m.search.IsMatch([]int{dbIdx, schemaIdx, tableIdx, colIdx})

	m.renderNode(t, colText, isFocused, isSearchMatch, isActiveMatch, nil)
}

// RenderContent returns the string representation of the view for composition
func (m DBTreeModel) RenderContent() string {
	t := tree.
		Root("Databases:").
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle()).
		RootStyle(rootStyle())

	// Render all databases
	for dbIdx, db := range m.tree.databases {
		m.renderDatabase(t, dbIdx, db)
	}

	fullContent := t.String()

	// TODO: This is probably not the best way to cut off content.
	// I think we shouldn't calculate everything and then cut lines,
	// but rather calculate only what fits in the window.
	lines := strings.Split(fullContent, "\n")
	if m.viewport.ScrollOffset() >= len(lines) {
		m.viewport.SetScrollOffset(max(0, len(lines)-m.viewport.Height()))
	}

	// Leave room for search bar if needed
	availableHeight := m.viewport.Height()
	if m.search.mode || len(m.search.matches) > 0 {
		availableHeight-- // Reserve one line for search bar
	}

	if len(lines) > availableHeight {
		end := min(m.viewport.ScrollOffset()+availableHeight, len(lines))
		lines = lines[m.viewport.ScrollOffset():end]
	}

	result := strings.Join(lines, "\n")

	// Apply width constraint
	resultLines := strings.Split(result, "\n")
	for i, line := range resultLines {
		if lipgloss.Width(line) > m.viewport.Width() {
			// Truncate lines that are too long
			resultLines[i] = line[:m.viewport.Width()]
		}
	}

	content := strings.Join(resultLines, "\n")

	// Add search bar if in search mode or has matches
	if m.search.mode || len(m.search.matches) > 0 {
		content += "\n" + m.renderSearchBar()
	}

	return content
}

// renderSearchBar renders the search input bar.
func (m DBTreeModel) renderSearchBar() string {
	style := searchBarStyle()

	if m.search.mode {
		// Active search input
		cursor := "_"
		matchInfo := ""
		if len(m.search.matches) > 0 {
			matchInfo = fmt.Sprintf(" [%d/%d]", m.search.matchIndex+1, len(m.search.matches))
		} else if m.search.query != "" {
			matchInfo = " [no matches]"
		}
		return style.Render(fmt.Sprintf("/%s%s%s", m.search.query, cursor, matchInfo))
	}

	// Not in search mode but showing match count
	if len(m.search.matches) > 0 {
		return style.Render(fmt.Sprintf("Search: %q [%d/%d] (n/N to navigate)",
			m.search.query, m.search.matchIndex+1, len(m.search.matches)))
	}

	return ""
}

// View implements tea.Model interface
func (m DBTreeModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}
