package tableview

type Viewport struct {
	width  int
	height int
}

// SetSize updates the dimensions of the TableView view
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

func (v *Viewport) VisibleHeight() int {
	return v.height - 3 // Reserve space for UI elements
}

func (v *Viewport) IsReady() bool {
	return v.width > 0 && v.height > 0
}
