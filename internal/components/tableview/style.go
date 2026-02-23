package tableview

import (
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// Status bar style functions

func sbIconStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Info).
		Background(theme.Current().Colors.Base).
		Bold(true)
}

func sbSepStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted).
		Background(theme.Current().Colors.Base)
}

func sbLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Subtle).
		Background(theme.Current().Colors.Base)
}

func sbValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Text).
		Background(theme.Current().Colors.Base)
}

func sbPaginationMsgStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Warning).
		Background(theme.Current().Colors.Base).
		Bold(true)
}

func sbInputLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Info).
		Background(theme.Current().Colors.Base).
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
		Foreground(theme.Current().Colors.Info).
		Background(theme.Current().Colors.Base)
}

func sbDimStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted).
		Background(theme.Current().Colors.Base)
}

// Table view style functions

func placeholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Muted).
		Background(theme.Current().Colors.Base).
		Italic(true)
}

func emptyStateCardStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Border).
		Background(colors.Surface).
		Padding(1, 2)
}

func emptyStateTitleStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Primary).
		Background(colors.Surface).
		Bold(true)
}

func emptyStateSubtitleStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Subtle).
		Background(colors.Surface)
}

func emptyStateHintStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Text).
		Background(colors.Surface)
}

func emptyStateBulletStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Info).
		Background(colors.Surface).
		Bold(true)
}

func spinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Colors.Primary)
}
