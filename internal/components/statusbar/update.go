package statusbar

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/pkgs/editor"
)

func (s StatusBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case editor.EditorCursorMovedMsg:
		cursorPos := fmt.Sprintf("%d:%d", msg.Row, msg.Col)
		s.editorCursorPos = cursorPos
	case editor.EditorModeChangedMsg:
		switch msg.Mode {
		case editor.EditorModeInsert:
			s.editorMode = "INSERT"
		case editor.EditorModeNormal:
			s.editorMode = "NORMAL"
		default:
			s.editorMode = "UNKNOWN"
		}
	}
	return s, tea.Batch(cmds...)
}

type UpdateStatusBarMsg struct {
	component string
	message   string
}
