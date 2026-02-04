package statusbar

import tea "charm.land/bubbletea/v2"

type StatusBarModel struct {
	width        int
	height       int
	editorStatus string
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		editorStatus: "Normal",
	}
}

func (s *StatusBarModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func UpdateStatusBar(component, message string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusBarMsg{
			component: component,
			message:   message,
		}
	}
}
