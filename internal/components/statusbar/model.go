package statusbar

type StatusBarModel struct {
	width  int
	height int
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{}
}

func (s *StatusBarModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}
