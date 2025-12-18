// Package theme provides a centralized theming system for the dbettier application.
// It defines color palettes and component styles that can be easily swapped
// to support different visual themes like Catppuccin, Rosé Pine, Gruvbox, etc.
package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme represents a complete color theme with semantic color roles
type Theme struct {
	Name   string
	Colors Colors
}

// Colors defines semantic color roles used throughout the application.
// These abstract away specific hex values into meaningful roles,
// making it easy to create consistent themes.
type Colors struct {
	// Backgrounds (darkest to lightest)
	Base    color.Color // Primary background - main panels, frames
	Surface color.Color // Secondary background - cards, inputs, sidebars
	Overlay color.Color // Tertiary background - popups, dialogs, tooltips

	// Foregrounds (lowest to highest contrast)
	Muted  color.Color // Lowest contrast - disabled, ignored content
	Subtle color.Color // Medium contrast - comments, secondary text
	Text   color.Color // Highest contrast - primary text, active content

	// Accent colors
	Primary   color.Color // Main accent - focus indicators, active elements
	Secondary color.Color // Secondary accent - links, highlights

	// Semantic colors
	Success color.Color // Green - success states, additions, confirmations
	Warning color.Color // Yellow/Orange - warnings, cautions
	Error   color.Color // Red - errors, deletions, critical
	Info    color.Color // Blue/Cyan - information, hints

	// UI-specific colors
	Border        color.Color // Default border color
	BorderFocused color.Color // Focused/active border color
	Selection     color.Color // Selection/highlight background
	SearchMatch   color.Color // Search match highlight
	SearchActive  color.Color // Active search match highlight

	// Extended palette (for syntax highlighting and special UI elements)
	Pink   color.Color // Functions, methods
	Peach  color.Color // Numbers, constants
	Yellow color.Color // Strings, warnings
	Green  color.Color // Success, additions
	Teal   color.Color // Types, special
	Blue   color.Color // Keywords, info
	Purple color.Color // Accent, special highlights
}

// current holds the active theme (default: Catppuccin Mocha)
var current = CatppuccinMocha()

// Current returns the currently active theme
func Current() Theme {
	return current
}

// SetTheme changes the active theme
func SetTheme(t Theme) {
	current = t
}

// AvailableThemes returns a list of all built-in themes
func AvailableThemes() []Theme {
	return []Theme{
		CatppuccinMocha(),
		CatppuccinMacchiato(),
		CatppuccinFrappe(),
		RosePine(),
		RosePineMoon(),
		GruvboxDark(),
		TokyoNight(),
		TokyoNightStorm(),
	}
}

// ThemeNames returns the names of all available themes
func ThemeNames() []string {
	themes := AvailableThemes()
	names := make([]string, len(themes))
	for i, t := range themes {
		names[i] = t.Name
	}
	return names
}

// GetThemeByName returns a theme by name, or the default theme if not found
func GetThemeByName(name string) Theme {
	for _, t := range AvailableThemes() {
		if t.Name == name {
			return t
		}
	}
	return CatppuccinMocha() // Default fallback
}

// Helper function to create color from hex string
func hex(s string) color.Color {
	return lipgloss.Color(s)
}
