package logpanel

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/messages"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// getStyleForLevel returns the lipgloss style for a given log level
func getStyleForLevel(level messages.LogLevel) lipgloss.Style {
	colors := theme.Current().Colors
	base := lipgloss.NewStyle().Background(colors.Base)
	switch level {
	case messages.LogInfo:
		return base.Foreground(colors.Subtle)
	case messages.LogSuccess:
		return base.Foreground(colors.Success)
	case messages.LogWarning:
		return base.Foreground(colors.Warning)
	case messages.LogError:
		return base.Foreground(colors.Error)
	case messages.LogSQL:
		return base.Foreground(colors.Info)
	default:
		return base
	}
}
