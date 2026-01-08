package workspace

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// TabBarHeight is the height of the tab bar (including borders)
const TabBarHeight = 3

// activeTabBorder creates a border for active tabs with open bottom
func activeTabBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
}

// inactiveTabBorder creates a border for inactive tabs
func inactiveTabBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
}

// tabGapBorder creates a border for the gap between tabs and edge
func tabGapBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: "",
	}
}

// activeTabStyle returns the style for the active tab
func activeTabStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Border(activeTabBorder(), true).
		BorderForeground(colors.Primary).
		Foreground(colors.Text).
		Padding(0, 1)
}

// inactiveTabStyle returns the style for inactive tabs
func inactiveTabStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Border(inactiveTabBorder(), true).
		BorderForeground(colors.Border).
		Foreground(colors.Subtle).
		Padding(0, 1)
}

// tabGapStyle returns the style for the gap/filler at the end of tabs
func tabGapStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		BorderStyle(tabGapBorder()).
		BorderBottom(true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(colors.Border)
}

// closeButtonStyle returns the style for the close button
func closeButtonStyle(active bool) lipgloss.Style {
	colors := theme.Current().Colors
	style := lipgloss.NewStyle().
		MarginLeft(1)
	if active {
		style = style.Foreground(colors.Error)
	} else {
		style = style.Foreground(colors.Muted)
	}
	return style
}

// iconStyle returns the style for tab icons
func iconStyle(tabType TabType) lipgloss.Style {
	colors := theme.Current().Colors
	style := lipgloss.NewStyle().MarginRight(1)
	switch tabType {
	case TabTypeTable:
		style = style.Foreground(colors.Blue)
	case TabTypeQuery:
		style = style.Foreground(colors.Purple)
	}
	return style
}
