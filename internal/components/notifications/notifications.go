// Package notifications provides toast-style notification components for user feedback.
package notifications

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type NotificationLevel int

const (
	InfoNotification NotificationLevel = iota
	SuccessNotification
	WarningNotification
	ErrorNotification
)

type Notification struct {
	Message string
	Level   NotificationLevel
}
type (
	ShowNotificationMsg  Notification
	ClearNotificationMsg struct{}
)

func (n *Notification) GetStyle() lipgloss.Style {
	switch n.Level {
	case SuccessNotification:
		return notificationSuccessStyle()
	case WarningNotification:
		return notificationWarningStyle()
	case ErrorNotification:
		return notificationErrorStyle()
	default:
		return notificationInfoStyle()
	}
}

func ShowSuccess(message string) tea.Cmd {
	return func() tea.Msg {
		return ShowNotificationMsg{
			Message: message,
			Level:   SuccessNotification,
		}
	}
}

func ShowError(message string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Error notification: %s", message)
		return ShowNotificationMsg{
			Message: message,
			Level:   ErrorNotification,
		}
	}
}

func ShowWarning(message string) tea.Cmd {
	return func() tea.Msg {
		return ShowNotificationMsg{
			Message: message,
			Level:   WarningNotification,
		}
	}
}

func ShowInfo(message string) tea.Cmd {
	return func() tea.Msg {
		return ShowNotificationMsg{
			Message: message,
			Level:   InfoNotification,
		}
	}
}
