package statusbar

import tea "charm.land/bubbletea/v2"

func (s StatusBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	return s, tea.Batch(cmds...)
}
