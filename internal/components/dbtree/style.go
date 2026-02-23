package dbtree

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// Style functions that use the current theme

func enumeratorStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Secondary).
		Background(colors.Base).
		MarginRight(1)
}

func rootStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Success).
		Background(colors.Base)
}

func itemStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Text).
		Background(colors.Base)
}

func focusedStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Primary).
		Background(colors.Base)
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
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Subtle).
		Background(colors.Base)
}
