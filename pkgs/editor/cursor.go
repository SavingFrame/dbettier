package editor

import (
	"image/color"
	"time"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// CursorStyle is the style for real and virtual cursors.
type CursorStyle struct {
	// Style styles the cursor block.
	//
	// For real cursors, the foreground color set here will be used as the
	// cursor color.
	Color color.Color

	// Shape is the cursor shape. The following shapes are available:
	//
	// - tea.CursorBlock
	// - tea.CursorUnderline
	// - tea.CursorBar
	//
	// This is only used for real cursors.
	Shape tea.CursorShape

	// CursorBlink determines whether or not the cursor should blink.
	Blink bool

	// BlinkSpeed is the speed at which the virtual cursor blinks. This has no
	// effect on real cursors as well as no effect if the cursor is set not to
	// [CursorBlink].
	//
	// By default, the blink speed is set to about 500ms.
	BlinkSpeed time.Duration
}

type editorCursor struct {
	row           int
	col           int
	virtualCursor cursor.Model
	style         CursorStyle
}

func newEditorCursor(row, col int) *editorCursor {
	return &editorCursor{
		row:           row,
		col:           col,
		virtualCursor: cursor.New(),
		style: CursorStyle{
			Color: lipgloss.Color("7"),
			Shape: tea.CursorBlock,
			Blink: true,
		},
	}
}

func (c *editorCursor) setPosition(row, col int) {
	c.row = row
	c.col = col
}

func (c *editorCursor) moveUp(n int, buff *buffer) {
	if c.row-n >= 0 {
		c.row -= n
	}

	if c.col > len(buff.lines[c.row]) {
		c.col = len(buff.lines[c.row])
	}
}

func (c *editorCursor) moveDown(n int, buff *buffer) {
	if c.row+n < len(buff.lines) {
		c.row += n
	} else {
		c.row = len(buff.lines) - 1
	}

	if c.col > len(buff.lines[c.row]) {
		c.col = len(buff.lines[c.row])
	}
}

func (c *editorCursor) moveLeft(n int) {
	if c.col-n >= 0 {
		c.col -= n
	}
}

func (c *editorCursor) moveRight(n int, buff *buffer) {
	line := buff.lines[c.row]
	if c.col+n <= len(line) {
		c.col += n
	} else {
		c.col = len(line)
	}
}

func (c *editorCursor) moveLastSymbol(lines []string) {
	c.col = len(lines[c.row]) - 1
	c.row = len(lines) - 1
}

func (c *editorCursor) gotoStartEdge(buff *buffer) {
	line := buff.lines[c.row]
	for i, char := range line {
		// NOTE: This is probably wrong for non-ASCII characters, but it should work for most cases.
		if char != ' ' {
			c.col = i
			break
		}
	}
}

func (c *editorCursor) gotoEndEdge(buff *buffer) {
	line := buff.lines[c.row]
	for i := len(line) - 1; i >= 0; i-- {
		// NOTE: This is probably wrong for non-ASCII characters, but it should work for most cases.
		if line[i] != ' ' {
			c.col = i
			break
		}
	}
}

func (c *editorCursor) setCharUnderCursor(lines []string) {
	focusedLine := lines[c.row]
	var charUnderCursor string
	if c.col < len(focusedLine) {
		charUnderCursor = string(focusedLine[c.col])
	} else {
		charUnderCursor = " "
	}
	c.virtualCursor.SetChar(charUnderCursor)
}
