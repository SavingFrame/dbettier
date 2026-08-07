package sqlcommandbar

import (
	"os"
	"os/exec"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/query"
)

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SQLCommandBarKeymap.Execute):
			q := m.editor.GetContent()
			return m, func() tea.Msg {
				return messages.ExecuteSQLTextMsg{
					Query:      q,
					DatabaseID: m.DatabaseID,
				}
			}
		case key.Matches(msg, SQLCommandBarKeymap.OpenNvim):
			name, err := createTempBufferFile(&m.editor)
			if err != nil {
				return m, notifications.ShowError("Failed to connect to database: " + err.Error())
			}
			cmd := exec.Command("nvim", name)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return messages.NvimFinishedMsg{Error: err, FileName: name}
			})
		}
	case query.SQLResultMsg:
		m.SetContent(msg.Query.Compile())
	case messages.NvimFinishedMsg:
		updateEditorFromFile(msg.FileName, &m.editor)
		deleteTempBufferFile(msg.FileName)
	}
	m.editor, cmd = m.editor.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
