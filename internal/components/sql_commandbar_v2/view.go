package sqlcommandbarv2

import tea "charm.land/bubbletea/v2"

func (m SQLCommandBarModel) RenderContent() string {
	return m.editor.View()
}

func (m SQLCommandBarModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}

func (m SQLCommandBarModel) Init() tea.Cmd {
	return nil
}
