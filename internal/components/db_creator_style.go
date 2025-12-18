package components

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

// DB Creator style functions using current theme

func dbcFocusedColor() color.Color {
	return theme.Current().Colors.Primary
}

func dbcFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Primary)
}

func dbcBlurredStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Muted)
}

func dbcSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Success)
}

func dbcErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Error)
}

func dbcHelpStyle() lipgloss.Style {
	return dbcBlurredStyle()
}

func dbcCursorModeHelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Colors.Subtle)
}

// Button rendering functions
func dbcFocusedButton() string {
	return dbcFocusedStyle().Render("[ Submit ]")
}

func dbcBlurredButton() string {
	return fmt.Sprintf("[ %s ]", dbcBlurredStyle().Render("Submit"))
}

func dbcFocusedTestButton() string {
	return dbcFocusedStyle().Render("[ Test Connection ]")
}

func dbcTestButton() string {
	return fmt.Sprintf("[ %s ]", dbcBlurredStyle().Render("Test Connection"))
}
