package logpanel

import (
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
)

func (m LogPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case sharedcomponents.AddLogMsg:
		m.AddLog(convertLogLevel(msg.Level), msg.Message)
		return m, nil
	case tea.KeyPressMsg, tea.MouseMsg, tea.MouseWheelMsg:
		// Forward to viewport for scrolling
		if m.ready {
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}

	return m, cmd
}

// convertLogLevel converts shared log level to local log level
func convertLogLevel(level sharedcomponents.LogLevel) sharedcomponents.LogLevel {
	switch level {
	case sharedcomponents.LogInfo:
		return sharedcomponents.LogInfo
	case sharedcomponents.LogSuccess:
		return sharedcomponents.LogSuccess
	case sharedcomponents.LogWarning:
		return sharedcomponents.LogWarning
	case sharedcomponents.LogError:
		return sharedcomponents.LogError
	case sharedcomponents.LogSQL:
		return sharedcomponents.LogSQL
	default:
		return sharedcomponents.LogInfo
	}
}
