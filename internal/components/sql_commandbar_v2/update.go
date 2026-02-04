package sqlcommandbarv2

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/components/statusbar"
	"github.com/SavingFrame/dbettier/pkgs/editor"
)

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, editor.NormalModeKeymap.EnableInsertMode) && m.editor.Mode() == editor.EditorModeNormal {
			cmd = statusbar.UpdateStatusBar("editorStatus", "INSERT")
			cmds = append(cmds, cmd)
		} else if key.Matches(msg, editor.InsertModeKeymap.Exit) && m.editor.Mode() == editor.EditorModeInsert {
			cmd = statusbar.UpdateStatusBar("editorStatus", "NORMAL")
			cmds = append(cmds, cmd)
		}

	case sharedcomponents.SQLResultMsg:
		m.editor.SetContent(msg.Query.Compile())
	}
	m.editor, cmd = m.editor.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
