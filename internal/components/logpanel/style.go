package logpanel

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// getStyleForLevel returns the lipgloss style for a given log level
func getStyleForLevel(level messages.LogLevel) lipgloss.Style {
	colors := theme.Current().Colors
	switch level {
	case messages.LogInfo:
		return lipgloss.NewStyle().Foreground(colors.Subtle)
	case messages.LogSuccess:
		return lipgloss.NewStyle().Foreground(colors.Success)
	case messages.LogWarning:
		return lipgloss.NewStyle().Foreground(colors.Warning)
	case messages.LogError:
		return lipgloss.NewStyle().Foreground(colors.Error)
	case messages.LogSQL:
		return lipgloss.NewStyle().Foreground(colors.Info)
	default:
		return lipgloss.NewStyle()
	}
}
