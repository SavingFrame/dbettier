package logpanel

import (
	tea "charm.land/bubbletea/v2"
)

// RenderContent returns the string representation of the view for composition
func (m LogPanelModel) RenderContent() string {
	if !m.ready {
		return "Initializing logs..."
	}
	return m.viewport.View()
}

// View implements tea.Model interface
func (m LogPanelModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.RenderContent())
	return v
}
