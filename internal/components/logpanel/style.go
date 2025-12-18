package logpanel

import (
	"charm.land/lipgloss/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// getStyleForLevel returns the lipgloss style for a given log level
func getStyleForLevel(level sharedcomponents.LogLevel) lipgloss.Style {
	colors := theme.Current().Colors
	switch level {
	case sharedcomponents.LogInfo:
		return lipgloss.NewStyle().Foreground(colors.Subtle)
	case sharedcomponents.LogSuccess:
		return lipgloss.NewStyle().Foreground(colors.Success)
	case sharedcomponents.LogWarning:
		return lipgloss.NewStyle().Foreground(colors.Warning)
	case sharedcomponents.LogError:
		return lipgloss.NewStyle().Foreground(colors.Error)
	case sharedcomponents.LogSQL:
		return lipgloss.NewStyle().Foreground(colors.Info)
	default:
		return lipgloss.NewStyle()
	}
}
