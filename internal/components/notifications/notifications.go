package notifications

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var (
	notficationInfoStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 2).
				Bold(true)

	notificationSuccessStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("42")).
					Foreground(lipgloss.Color("0")).
					Padding(0, 2).
					Bold(true)

	notificationWarningStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("214")).
					Foreground(lipgloss.Color("0")).
					Padding(0, 2).
					Bold(true)

	notificationErrorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("196")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 2).
				Bold(true)
)

func (n *Notification) GetStyle() lipgloss.Style {
	switch n.Level {
	case SuccessNotification:
		return notificationSuccessStyle
	case WarningNotification:
		return notificationWarningStyle
	case ErrorNotification:
		return notificationErrorStyle
	default:
		return notficationInfoStyle
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
