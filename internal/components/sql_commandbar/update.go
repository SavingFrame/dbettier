package sqlcommandbar

import (
	"log"

	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

func (m SQLCommandBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case sharedcomponents.SetSQLTextMsg:
		log.Println("Get `SETSQLTEXT` message in SQLCommandBarModel")
		m.textarea.SetValue(msg.Command)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case "ctrl+c":
			return m, tea.Quit
			// Don't auto-focus - let root handle focus
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Focus sets focus to the textarea
func (m *SQLCommandBarModel) Focus() tea.Cmd {
	return m.textarea.Focus()
}

// Blur removes focus from the textarea
func (m *SQLCommandBarModel) Blur() {
	m.textarea.Blur()
}
