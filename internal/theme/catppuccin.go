package theme

// Catppuccin color palette
// https://github.com/catppuccin/catppuccin
//
// Catppuccin is a community-driven pastel theme that aims to be the middle ground
// between low and high contrast themes. It consists of 4 soothing warm flavors.

// CatppuccinMocha returns the Catppuccin Mocha theme (darkest flavor)
func CatppuccinMocha() Theme {
	return Theme{
		Name: "Catppuccin Mocha",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#1e1e2e"), // Base
			Surface: hex("#313244"), // Surface0
			Overlay: hex("#45475a"), // Surface1

			// Foregrounds
			Muted:  hex("#6c7086"), // Overlay0
			Subtle: hex("#a6adc8"), // Subtext0
			Text:   hex("#cdd6f4"), // Text

			// Accents
			Primary:   hex("#cba6f7"), // Mauve
			Secondary: hex("#89b4fa"), // Blue

			// Semantic
			Success: hex("#a6e3a1"), // Green
			Warning: hex("#f9e2af"), // Yellow
			Error:   hex("#f38ba8"), // Red
			Info:    hex("#89dceb"), // Sky

			// UI specific
			Border:        hex("#45475a"), // Surface1
			BorderFocused: hex("#f5c2e7"), // Pink
			Selection:     hex("#585b70"), // Surface2
			SearchMatch:   hex("#f9e2af"), // Yellow
			SearchActive:  hex("#fab387"), // Peach

			// Extended palette
			Pink:   hex("#f5c2e7"), // Pink
			Peach:  hex("#fab387"), // Peach
			Yellow: hex("#f9e2af"), // Yellow
			Green:  hex("#a6e3a1"), // Green
			Teal:   hex("#94e2d5"), // Teal
			Blue:   hex("#89b4fa"), // Blue
			Purple: hex("#cba6f7"), // Mauve
		},
	}
}

// CatppuccinMacchiato returns the Catppuccin Macchiato theme (medium-dark flavor)
func CatppuccinMacchiato() Theme {
	return Theme{
		Name: "Catppuccin Macchiato",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#24273a"), // Base
			Surface: hex("#363a4f"), // Surface0
			Overlay: hex("#494d64"), // Surface1

			// Foregrounds
			Muted:  hex("#6e738d"), // Overlay0
			Subtle: hex("#a5adcb"), // Subtext0
			Text:   hex("#cad3f5"), // Text

			// Accents
			Primary:   hex("#c6a0f6"), // Mauve
			Secondary: hex("#8aadf4"), // Blue

			// Semantic
			Success: hex("#a6da95"), // Green
			Warning: hex("#eed49f"), // Yellow
			Error:   hex("#ed8796"), // Red
			Info:    hex("#91d7e3"), // Sky

			// UI specific
			Border:        hex("#494d64"), // Surface1
			BorderFocused: hex("#f5bde6"), // Pink
			Selection:     hex("#5b6078"), // Surface2
			SearchMatch:   hex("#eed49f"), // Yellow
			SearchActive:  hex("#f5a97f"), // Peach

			// Extended palette
			Pink:   hex("#f5bde6"), // Pink
			Peach:  hex("#f5a97f"), // Peach
			Yellow: hex("#eed49f"), // Yellow
			Green:  hex("#a6da95"), // Green
			Teal:   hex("#8bd5ca"), // Teal
			Blue:   hex("#8aadf4"), // Blue
			Purple: hex("#c6a0f6"), // Mauve
		},
	}
}

// CatppuccinFrappe returns the Catppuccin Frappé theme (medium flavor)
func CatppuccinFrappe() Theme {
	return Theme{
		Name: "Catppuccin Frappe",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#303446"), // Base
			Surface: hex("#414559"), // Surface0
			Overlay: hex("#51576d"), // Surface1

			// Foregrounds
			Muted:  hex("#737994"), // Overlay0
			Subtle: hex("#a5adce"), // Subtext0
			Text:   hex("#c6d0f5"), // Text

			// Accents
			Primary:   hex("#ca9ee6"), // Mauve
			Secondary: hex("#8caaee"), // Blue

			// Semantic
			Success: hex("#a6d189"), // Green
			Warning: hex("#e5c890"), // Yellow
			Error:   hex("#e78284"), // Red
			Info:    hex("#99d1db"), // Sky

			// UI specific
			Border:        hex("#51576d"), // Surface1
			BorderFocused: hex("#f4b8e4"), // Pink
			Selection:     hex("#626880"), // Surface2
			SearchMatch:   hex("#e5c890"), // Yellow
			SearchActive:  hex("#ef9f76"), // Peach

			// Extended palette
			Pink:   hex("#f4b8e4"), // Pink
			Peach:  hex("#ef9f76"), // Peach
			Yellow: hex("#e5c890"), // Yellow
			Green:  hex("#a6d189"), // Green
			Teal:   hex("#81c8be"), // Teal
			Blue:   hex("#8caaee"), // Blue
			Purple: hex("#ca9ee6"), // Mauve
		},
	}
}
