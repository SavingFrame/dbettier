package table

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	colors := theme.Current().Colors
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colors.Border),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		SelectedCell: lipgloss.NewStyle().
			Padding(0, 1).
			Background(colors.Primary).
			Foreground(colors.Base).
			Bold(true),
		SelectedRow: lipgloss.NewStyle().
			Padding(0, 1).
			Background(colors.Surface),
		SelectedCol: lipgloss.NewStyle().
			Padding(0, 1).
			Background(colors.Surface),
		SearchMatch: lipgloss.NewStyle().
			Padding(0, 1).
			Background(colors.SearchMatch).
			Foreground(colors.Base),
		SearchMatchActive: lipgloss.NewStyle().
			Padding(0, 1).
			Background(colors.SearchActive).
			Foreground(colors.Base).
			Bold(true),
	}
}

// headerFocusedStyle returns style for a focused column header.
func headerFocusedStyle(base lipgloss.Style) lipgloss.Style {
	return base.Background(theme.Current().Colors.Primary).Bold(true)
}

// searchBarStyle returns the style for the search bar.
func searchBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Subtle)
}

// scrollIndicatorStyle returns the style for scroll indicators.
func scrollIndicatorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Subtle)
}
