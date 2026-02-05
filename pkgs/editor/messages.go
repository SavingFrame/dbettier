package editor

// EditorModeChangedMsg is emitted when the editor mode changes
type EditorModeChangedMsg struct {
	Mode EditorMode
}

// EditorCursorMovedMsg is emitted when the cursor position changes
type EditorCursorMovedMsg struct {
	Row int
	Col int
}
