package sqlcommandbarv2

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/query"
)

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SQLCommandBarV2Keymap.Execute):
			q := m.editor.GetContent()
			return m, func() tea.Msg {
				return messages.ExecuteSQLTextMsg{
					Query:      q,
					DatabaseID: m.DatabaseID,
				}
			}
		}
	case query.SQLResultMsg:
		m.SetContent(msg.Query.Compile())
	}
	m.editor, cmd = m.editor.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
