// Package logpanel provides a scrollable log panel component for displaying application logs.
package logpanel

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
)

const MaxLogEntries = 1000

// LogEntry represents a single log message with styling
type LogEntry struct {
	Message string
	Level   sharedcomponents.LogLevel
}

// LogPanelModel handles the log display viewport
type LogPanelModel struct {
	viewport viewport.Model
	entries  []LogEntry
	width    int
	height   int
	ready    bool
}

// LogPanelScreen creates a new log panel
func LogPanelScreen() LogPanelModel {
	return LogPanelModel{
		entries: make([]LogEntry, 0, MaxLogEntries),
	}
}

func (m LogPanelModel) Init() tea.Cmd {
	return nil
}

// SetSize updates the dimensions of the log panel
func (m *LogPanelModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	if !m.ready {
		m.viewport = viewport.New(
			viewport.WithWidth(width),
			viewport.WithHeight(height),
		)
		m.viewport.LeftGutterFunc = func(info viewport.GutterContext) string {
			if info.Soft {
				return "   "
			}
			return ""
		}
		m.ready = true
	} else {
		m.viewport.SetWidth(width)
		m.viewport.SetHeight(height)
	}
	m.refreshContent()
}

// AddEntry adds a new log entry and maintains the buffer limit
func (m *LogPanelModel) AddEntry(entry LogEntry) {
	m.entries = append(m.entries, entry)

	// Trim old entries if we exceed the limit
	if len(m.entries) > MaxLogEntries {
		m.entries = m.entries[len(m.entries)-MaxLogEntries:]
	}

	m.refreshContent()
	// Auto-scroll to bottom
	m.viewport.GotoBottom()
}

// AddLog is a convenience method to add a log with level
func (m *LogPanelModel) AddLog(level sharedcomponents.LogLevel, message string) {
	m.AddEntry(LogEntry{
		Message: message,
		Level:   level,
	})
}

// Clear removes all log entries
func (m *LogPanelModel) Clear() {
	m.entries = m.entries[:0]
	m.refreshContent()
}

// refreshContent rebuilds the viewport content from entries
func (m *LogPanelModel) refreshContent() {
	if !m.ready {
		return
	}

	var content string
	for i, entry := range m.entries {
		style := getStyleForLevel(entry.Level)
		line := style.Render(entry.Message)
		if i > 0 {
			content += "\n"
		}
		content += line
	}

	m.viewport.SetContent(content)
}

func AddLogCmd(message string, level sharedcomponents.LogLevel) tea.Cmd {
	return func() tea.Msg {
		return sharedcomponents.AddLogMsg{
			Message: message,
			Level:   level,
		}
	}
}
