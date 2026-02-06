package logpanel

import (
	"bytes"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/theme"
	"github.com/alecthomas/chroma/v2/quick"
)

func (m LogPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.AddLogMsg:
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
func convertLogLevel(level messages.LogLevel) messages.LogLevel {
	switch level {
	case messages.LogInfo:
		return messages.LogInfo
	case messages.LogSuccess:
		return messages.LogSuccess
	case messages.LogWarning:
		return messages.LogWarning
	case messages.LogError:
		return messages.LogError
	case messages.LogSQL:
		return messages.LogSQL
	default:
		return messages.LogInfo
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
