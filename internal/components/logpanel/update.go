package logpanel

import (
	"bytes"
	"fmt"

	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/theme"
	"github.com/alecthomas/chroma/v2/quick"
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

// highlightCode returns a syntax highlighted string of text.
func highlightCode(content, extension string) (string, error) {
	buf := new(bytes.Buffer)
	s := theme.Current().Name
	if err := quick.Highlight(buf, content, extension, "terminal256", s); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return buf.String(), nil
}
