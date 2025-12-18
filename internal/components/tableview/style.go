package tableview

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// Status bar style functions

func sbIconStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Info).
		Bold(true)
}

func sbSepStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted)
}

func sbLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Subtle)
}

func sbValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Text)
}

func sbPaginationMsgStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Warning).
		Bold(true)
}

func sbInputLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Info).
		Bold(true)
}

func sbInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Current().Colors.Surface).
		Padding(0, 1)
}

func sbButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Text).
		Background(theme.Current().Colors.Overlay).
		Padding(0, 1)
}

func sbButtonPrimaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Base).
		Background(theme.Current().Colors.Primary).
		Bold(true).
		Padding(0, 1)
}

func sbSortStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Info)
}

func sbDimStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted)
}

// Table view style functions

func placeholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted).
		Italic(true)
}

func spinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Primary)
}
