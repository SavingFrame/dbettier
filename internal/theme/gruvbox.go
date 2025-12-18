package theme

// Gruvbox color palette
// https://github.com/morhetz/gruvbox
//
// Gruvbox is a retro groove color scheme designed to be easy on the eyes.
// It provides both dark and light variants with warm, earthy colors.

// GruvboxDark returns the Gruvbox Dark theme (hard contrast)
func GruvboxDark() Theme {
	return Theme{
		Name: "Gruvbox Dark",
		Colors: Colors{
			// Backgrounds
			Base:    hex("#1d2021"), // bg0_h (hard contrast)
			Surface: hex("#282828"), // bg0
			Overlay: hex("#3c3836"), // bg1

			// Foregrounds
			Muted:  hex("#665c54"), // bg4
			Subtle: hex("#a89984"), // gray
			Text:   hex("#ebdbb2"), // fg

			// Accents
			Primary:   hex("#d3869b"), // Purple
			Secondary: hex("#83a598"), // Aqua

			// Semantic
			Success: hex("#b8bb26"), // Green
			Warning: hex("#fabd2f"), // Yellow
			Error:   hex("#fb4934"), // Red
			Info:    hex("#83a598"), // Aqua

			// UI specific
			Border:        hex("#504945"), // bg2
			BorderFocused: hex("#d3869b"), // Purple
			Selection:     hex("#504945"), // bg2
			SearchMatch:   hex("#fabd2f"), // Yellow
			SearchActive:  hex("#fe8019"), // Orange

			// Extended palette
			Pink:   hex("#d3869b"), // Purple
			Peach:  hex("#fe8019"), // Orange
			Yellow: hex("#fabd2f"), // Yellow
			Green:  hex("#b8bb26"), // Green
			Teal:   hex("#8ec07c"), // Aqua bright
			Blue:   hex("#83a598"), // Blue
			Purple: hex("#d3869b"), // Purple
		},
	}
}
