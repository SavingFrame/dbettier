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
		Background(colors.Base).
		Padding(0, 1)
}

// inactiveTabStyle returns the style for inactive tabs
func inactiveTabStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Border(inactiveTabBorder(), true).
		BorderForeground(colors.Border).
		Foreground(colors.Subtle).
		Background(colors.Surface).
		Padding(0, 1)
}

// tabGapStyle returns the style for the gap/filler at the end of tabs
func tabGapStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Background(colors.Base)
}

func emptyTabBarStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Background(colors.Base)
}

// closeButtonStyle returns the style for the close button
func closeButtonStyle(active bool) lipgloss.Style {
	colors := theme.Current().Colors
	style := lipgloss.NewStyle().
		MarginLeft(1)
	if active {
		style = style.Foreground(colors.Error).Background(colors.Base)
	} else {
		style = style.Foreground(colors.Muted).Background(colors.Surface)
	}
	return style
}

// iconStyle returns the style for tab icons
func iconStyle(tabType TabType, active bool) lipgloss.Style {
	colors := theme.Current().Colors
	bg := colors.Surface
	if active {
		bg = colors.Base
	}

	style := lipgloss.NewStyle().
		MarginRight(1).
		Background(bg)
	switch tabType {
	case TabTypeTable:
		style = style.Foreground(colors.Blue)
	case TabTypeQuery:
		style = style.Foreground(colors.Purple)
	}
	return style
}

func tabNameStyle(active bool) lipgloss.Style {
	colors := theme.Current().Colors
	fg := colors.Subtle
	bg := colors.Surface
	if active {
		fg = colors.Text
		bg = colors.Base
	}
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg)
}

func tabSpaceStyle(active bool) lipgloss.Style {
	colors := theme.Current().Colors
	bg := colors.Surface
	if active {
		bg = colors.Base
	}
	return lipgloss.NewStyle().Background(bg)
}

func emptyWorkspaceCardStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Border).
		Background(colors.Surface).
		Padding(1, 2)
}

func emptyWorkspaceTitleStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Primary).
		Background(colors.Surface).
		Bold(true)
}

func emptyWorkspaceSubtitleStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Subtle).
		Background(colors.Surface)
}

func emptyWorkspaceHintStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Text).
		Background(colors.Surface)
}

func emptyWorkspaceBulletStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Info).
		Background(colors.Surface).
		Bold(true)
}
