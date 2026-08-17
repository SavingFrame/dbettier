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
			dir, err := createTempDir()
			if err != nil {
				return m, notifications.ShowError("Failed to create temperary directory: " + err.Error())
			}
			queryFile, err := createTempContentFile(dir, &m.editor)
			if err != nil {
				return m, notifications.ShowError("Failed to create temperary content file: " + err.Error())
			}
			_, err = createTempCredentialsFile(dir, *m.registry.GetByID(m.DatabaseID))
			if err != nil {
				return m, notifications.ShowError("Failed to create temperary credentails file: " + err.Error())
			}
			cmd := exec.Command("nvim", "--cmd", "lua vim.opt.runtimepath:prepend(vim.env.DBETTIER_NVIM_RUNTIME)", queryFile)
			runtimeDir := "~/Projects/dbettier/resources/nvim"
			cmd.Env = append(os.Environ(),
				"DBETTIER_NVIM_RUNTIME="+runtimeDir,
				"DBETTIER_PGLS_ROOT="+dir,
			)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return messages.NvimFinishedMsg{Error: err, FileName: queryFile}
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
