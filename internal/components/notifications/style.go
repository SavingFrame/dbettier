package notifications

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// Notification style functions using current theme

func notificationInfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Current().Colors.Info).
		Foreground(theme.Current().Colors.Base).
		Padding(0, 2).
		Bold(true)
}

func notificationSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Current().Colors.Success).
		Foreground(theme.Current().Colors.Base).
		Padding(0, 2).
		Bold(true)
}

func notificationWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Current().Colors.Warning).
		Foreground(theme.Current().Colors.Base).
		Padding(0, 2).
		Bold(true)
}

func notificationErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Current().Colors.Error).
		Foreground(theme.Current().Colors.Text).
		Padding(0, 2).
		Bold(true)
}
