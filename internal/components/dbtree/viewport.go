package dbtree

type Viewport struct {
	width        int
	height       int
	scrollOffset int
}

// SetSize updates the dimensions of the DBTree view
func (v *Viewport) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *Viewport) Width() int {
	return v.width
}

func (v *Viewport) Height() int {
	return v.height
}

func (v *Viewport) ScrollOffset() int {
	return v.scrollOffset
}

func (v *Viewport) SetScrollOffset(offset int) {
	v.scrollOffset = offset
	if v.scrollOffset < 0 {
		v.scrollOffset = 0
	}
}

func (v *Viewport) VisibleHeight() int {
	return v.height - 3 // Reserve space for UI elements
}

func (v *Viewport) AdjustScrollToCursor(cursorLine int) {
	visibleHeight := v.VisibleHeight()

	// If cursor is above viewport, scroll up
	if cursorLine < v.scrollOffset {
		v.scrollOffset = cursorLine
	}

	// If cursor is below viewport, scroll down
	if cursorLine >= v.scrollOffset+visibleHeight {
		v.scrollOffset = cursorLine - visibleHeight + 1
	}

	// Ensure scrollOffset doesn't go negative
	if v.scrollOffset < 0 {
		v.scrollOffset = 0
	}
}

