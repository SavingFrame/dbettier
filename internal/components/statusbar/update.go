package statusbar

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func (s StatusBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case UpdateStatusBarMsg:
		cmd := s.handleUpdateStatusBar(msg)
		cmds = append(cmds, cmd)
	}
	return s, tea.Batch(cmds...)
}

type UpdateStatusBarMsg struct {
	component string
	message   string
}

func (s *StatusBarModel) handleUpdateStatusBar(msg UpdateStatusBarMsg) tea.Cmd {
	log.Printf("Updating status bar: component=%s, message=%s", msg.component, msg.message)
	switch msg.component {
	case "editorStatus":
		s.editorStatus = msg.message
	}
	return nil
}
