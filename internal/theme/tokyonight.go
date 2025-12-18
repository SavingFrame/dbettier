package theme

// Tokyo Night color palette
// https://github.com/folke/tokyonight.nvim
//
// A clean, dark Neovim theme ported from the Visual Studio Code TokyoNight theme.
// It provides multiple variants: Night (darker), Storm (default), Day (light), and Moon.

// TokyoNight returns the Tokyo Night theme (darker variant)
func TokyoNight() Theme {
	return Theme{
		Name: "Tokyo Night",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#1a1b26"), // bg
			Surface: hex("#24283b"), // bg_dark
			Overlay: hex("#414868"), // bg_highlight

			// Foregrounds
			Muted:  hex("#565f89"), // comment
			Subtle: hex("#a9b1d6"), // fg_dark
			Text:   hex("#c0caf5"), // fg

			// Accents
			Primary:   hex("#bb9af7"), // Purple
			Secondary: hex("#7aa2f7"), // Blue

			// Semantic
			Success: hex("#9ece6a"), // Green
			Warning: hex("#e0af68"), // Yellow
			Error:   hex("#f7768e"), // Red
			Info:    hex("#7dcfff"), // Cyan

			// UI specific
			Border:        hex("#414868"), // bg_highlight
			BorderFocused: hex("#bb9af7"), // Purple
			Selection:     hex("#364a82"), // bg_visual
			SearchMatch:   hex("#e0af68"), // Yellow
			SearchActive:  hex("#ff9e64"), // Orange

			// Extended palette
			Pink:   hex("#ff007c"), // Magenta
			Peach:  hex("#ff9e64"), // Orange
			Yellow: hex("#e0af68"), // Yellow
			Green:  hex("#9ece6a"), // Green
			Teal:   hex("#73daca"), // Teal
			Blue:   hex("#7aa2f7"), // Blue
			Purple: hex("#bb9af7"), // Purple
		},
	}
}

// TokyoNightStorm returns the Tokyo Night Storm theme (default variant)
func TokyoNightStorm() Theme {
	return Theme{
		Name: "Tokyo Night Storm",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#24283b"), // bg
			Surface: hex("#1f2335"), // bg_dark
			Overlay: hex("#414868"), // bg_highlight

			// Foregrounds
			Muted:  hex("#565f89"), // comment
			Subtle: hex("#a9b1d6"), // fg_dark
			Text:   hex("#c0caf5"), // fg

			// Accents
			Primary:   hex("#bb9af7"), // Purple
			Secondary: hex("#7aa2f7"), // Blue

			// Semantic
			Success: hex("#9ece6a"), // Green
			Warning: hex("#e0af68"), // Yellow
			Error:   hex("#f7768e"), // Red
			Info:    hex("#7dcfff"), // Cyan

			// UI specific
			Border:        hex("#414868"), // bg_highlight
			BorderFocused: hex("#bb9af7"), // Purple
			Selection:     hex("#364a82"), // bg_visual
			SearchMatch:   hex("#e0af68"), // Yellow
			SearchActive:  hex("#ff9e64"), // Orange

			// Extended palette
			Pink:   hex("#ff007c"), // Magenta
			Peach:  hex("#ff9e64"), // Orange
			Yellow: hex("#e0af68"), // Yellow
			Green:  hex("#9ece6a"), // Green
			Teal:   hex("#73daca"), // Teal
			Blue:   hex("#7aa2f7"), // Blue
			Purple: hex("#bb9af7"), // Purple
		},
	}
}
