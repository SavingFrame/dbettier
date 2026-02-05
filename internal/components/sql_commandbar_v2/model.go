package sqlcommandbarv2

import (
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/SavingFrame/dbettier/pkgs/editor"
)

type SQLCommandBarModel struct {
	editor editor.SQLEditor

	registry   *database.DBRegistry
	DatabaseID string
}

func NewSQLCommandBarModel(lines []string, registry *database.DBRegistry, databaseID string, readonly bool) SQLCommandBarModel {
	return SQLCommandBarModel{
		editor:     editor.NewEditorModel(lines, readonly),
		registry:   registry,
		DatabaseID: databaseID,
	}
}

func (m *SQLCommandBarModel) SetSize(width, height int) {
	m.editor.SetSize(width, height)
}

func (m *SQLCommandBarModel) Focus() tea.Cmd {
	return m.editor.Focus()
}

func (m *SQLCommandBarModel) Blur() tea.Cmd {
	return m.editor.Blur()
}

func (m *SQLCommandBarModel) SetContent(content string) {
	m.editor.SetContent(content)
}
