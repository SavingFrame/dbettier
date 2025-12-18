package theme

// Rosé Pine color palette
// https://rosepinetheme.com/palette/ingredients
//
// All natural pine, faux fur and a bit of soho vibes for the classy minimalist.
// Rosé Pine has three variants: Main, Moon, and Dawn (light).

// RosePine returns the Rosé Pine Main theme (darkest variant)
func RosePine() Theme {
	return Theme{
		Name: "Rose Pine",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#191724"), // Base
			Surface: hex("#1f1d2e"), // Surface
			Overlay: hex("#26233a"), // Overlay

			// Foregrounds
			Muted:  hex("#6e6a86"), // Muted
			Subtle: hex("#908caa"), // Subtle
			Text:   hex("#e0def4"), // Text

			// Accents
			Primary:   hex("#c4a7e7"), // Iris
			Secondary: hex("#9ccfd8"), // Foam

			// Semantic
			Success: hex("#31748f"), // Pine
			Warning: hex("#f6c177"), // Gold
			Error:   hex("#eb6f92"), // Love
			Info:    hex("#9ccfd8"), // Foam

			// UI specific
			Border:        hex("#524f67"), // Highlight High
			BorderFocused: hex("#ebbcba"), // Rose
			Selection:     hex("#403d52"), // Highlight Med
			SearchMatch:   hex("#f6c177"), // Gold
			SearchActive:  hex("#ebbcba"), // Rose

			// Extended palette
			Pink:   hex("#ebbcba"), // Rose
			Peach:  hex("#f6c177"), // Gold
			Yellow: hex("#f6c177"), // Gold
			Green:  hex("#31748f"), // Pine
			Teal:   hex("#9ccfd8"), // Foam
			Blue:   hex("#9ccfd8"), // Foam
			Purple: hex("#c4a7e7"), // Iris
		},
	}
}

// RosePineMoon returns the Rosé Pine Moon theme (medium-dark variant)
func RosePineMoon() Theme {
	return Theme{
		Name: "Rose Pine Moon",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#232136"), // Base
			Surface: hex("#2a273f"), // Surface
			Overlay: hex("#393552"), // Overlay

			// Foregrounds
			Muted:  hex("#6e6a86"), // Muted
			Subtle: hex("#908caa"), // Subtle
			Text:   hex("#e0def4"), // Text

			// Accents
			Primary:   hex("#c4a7e7"), // Iris
			Secondary: hex("#9ccfd8"), // Foam

			// Semantic
			Success: hex("#3e8fb0"), // Pine
			Warning: hex("#f6c177"), // Gold
			Error:   hex("#eb6f92"), // Love
			Info:    hex("#9ccfd8"), // Foam

			// UI specific
			Border:        hex("#56526e"), // Highlight High
			BorderFocused: hex("#ea9a97"), // Rose
			Selection:     hex("#44415a"), // Highlight Med
			SearchMatch:   hex("#f6c177"), // Gold
			SearchActive:  hex("#ea9a97"), // Rose

			// Extended palette
			Pink:   hex("#ea9a97"), // Rose
			Peach:  hex("#f6c177"), // Gold
			Yellow: hex("#f6c177"), // Gold
			Green:  hex("#3e8fb0"), // Pine
			Teal:   hex("#9ccfd8"), // Foam
			Blue:   hex("#9ccfd8"), // Foam
			Purple: hex("#c4a7e7"), // Iris
		},
	}
}
