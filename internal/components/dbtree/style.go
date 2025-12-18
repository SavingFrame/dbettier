package dbtree

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// Style functions that use the current theme

func enumeratorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Secondary).MarginRight(1)
}

func rootStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Success)
}

func itemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Text)
}

func focusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Primary)
}

func searchMatchStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Base).
		Background(theme.Current().Colors.SearchMatch)
}

func searchMatchActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Base).
		Background(theme.Current().Colors.SearchActive)
}

func searchBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Subtle)
}
