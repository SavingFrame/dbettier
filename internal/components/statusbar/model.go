package statusbar

type StatusBarModel struct {
	width           int
	height          int
	editorMode      string
	editorCursorPos string
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		editorMode:      "NORMAL",
		editorCursorPos: "1:1",
	}
}

func (s *StatusBarModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}
