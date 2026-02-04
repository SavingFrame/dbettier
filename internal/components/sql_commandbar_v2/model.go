package sqlcommandbarv2

import (
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/SavingFrame/dbettier/pkgs/editor"
)

type SQLCommandBarModel struct {
	editor editor.SQLEditor

	registry *database.DBRegistry
}

func NewSQLCommandBarModel(lines []string, registry *database.DBRegistry) SQLCommandBarModel {
	return SQLCommandBarModel{
		editor:   editor.NewEditorModel(lines),
		registry: registry,
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
