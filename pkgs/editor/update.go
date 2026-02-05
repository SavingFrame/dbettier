package editor

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m SQLEditor) Update(msg tea.Msg) (SQLEditor, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == EditorModeNormal {
			cmd = m.processNormalModeKey(msg)
			cmds = append(cmds, cmd)
		} else if m.mode == EditorModeInsert {
			cmd = m.processInsertModeKey(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.cursor.virtualCursor, cmd = m.cursor.virtualCursor.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SQLEditor) processNormalModeKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, NormalModeKeymap.Left):
		m.cursor.moveLeft(1)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, NormalModeKeymap.Right):
		m.cursor.moveRight(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, NormalModeKeymap.Up):
		m.cursor.moveUp(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, NormalModeKeymap.Down):
		m.cursor.moveDown(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, NormalModeKeymap.EnableInsertMode) && !m.readonly:
		m.mode = EditorModeInsert
		cmd = func() tea.Msg { return EditorModeChangedMsg{Mode: m.mode} }
	case key.Matches(msg, NormalModeKeymap.InsertNextChar) && !m.readonly:
		m.cursor.moveRight(1, m.buffer)
		m.mode = EditorModeInsert
		moveCmd := func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
		modeCmd := func() tea.Msg { return EditorModeChangedMsg{Mode: m.mode} }
		cmd = tea.Batch(moveCmd, modeCmd)
	case key.Matches(msg, NormalModeKeymap.InsertNewLine) && !m.readonly:
		if m.cursor.row+1 == len(m.buffer.lines) {
			m.buffer.lines = append(m.buffer.lines, "")
		}
		m.cursor.setPosition(m.cursor.row+1, 0)
		m.mode = EditorModeInsert
		moveCmd := func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
		modeCmd := func() tea.Msg { return EditorModeChangedMsg{Mode: m.mode} }
		cmd = tea.Batch(moveCmd, modeCmd)
	case key.Matches(msg, NormalModeKeymap.Exit):
		cmd = tea.Quit
	}
	return cmd
}

func (m *SQLEditor) processInsertModeKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, InsertModeKeymap.Exit):
		m.mode = EditorModeNormal
		cmd = func() tea.Msg { return EditorModeChangedMsg{Mode: m.mode} }
	case key.Matches(msg, InsertModeKeymap.Left):
		m.cursor.moveLeft(1)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, InsertModeKeymap.Right):
		m.cursor.moveRight(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, InsertModeKeymap.Up):
		m.cursor.moveUp(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, InsertModeKeymap.Down):
		m.cursor.moveDown(1, m.buffer)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, InsertModeKeymap.Backspace) && !m.IsReadonly():
		m.buffer.handleBackspace(m.cursor)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case key.Matches(msg, InsertModeKeymap.Space) && !m.IsReadonly():
		m.buffer.handleSpace(m.cursor)
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	case len(msg.String()) == 1 && !m.IsReadonly():
		m.buffer.handleCharacterInput(m.cursor, msg.String())
		cmd = func() tea.Msg { return EditorCursorMovedMsg{Row: m.cursor.row, Col: m.cursor.col} }
	}
	return cmd
}
