package editor

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func (m SQLEditor) View() string {
	// The highlighter produces ANSI-colored text, so we must slice and measure
	// by visible cell width (ANSI escape sequences take no space). The x/ansi
	// helpers handle this safely for us.
	//
	// Rendering pipeline:
	// 1) Highlight the raw buffer text (no placeholders; avoid breaking lexing).
	// 2) Split into ANSI-colored lines.
	// 3) On the cursor line, cut [left][cell][right] using ANSI-aware slicing.
	// 4) Apply reverse video to the cell when the cursor is visible, preserving
	//    the cell's syntax colors.
	// 5) Reassemble and render through the viewport.
	//
	// Note: ASCII-only cursor columns for now; use grapheme-aware mapping later
	// if you need multibyte/emoji support.

	// Keep the virtual cursor in sync with the character under the cursor.
	m.cursor.setCharUnderCursor(m.buffer.lines)

	// Highlight the raw buffer text (no placeholders) so tokenization isn't disturbed.
	content := strings.Join(m.buffer.lines, "\n")
	highlighted := highlightCode(content)
	lines := strings.Split(highlighted, "\n")

	// Overlay the cursor on the highlighted line using ANSI-aware slicing.
	if m.cursor.row >= 0 && m.cursor.row < len(lines) {
		line := lines[m.cursor.row]
		cursorCol := m.cursor.col
		lineWidth := ansi.StringWidth(line)
		if cursorCol < 0 {
			cursorCol = 0
		}
		if cursorCol > lineWidth {
			cursorCol = lineWidth
		}

		// Cut preserves ANSI sequences so we don't break colors.
		left := ansi.Cut(line, 0, cursorCol)
		cell := " "
		if cursorCol < lineWidth {
			cell = ansi.Cut(line, cursorCol, cursorCol+1)
		}
		right := ""
		if cursorCol < lineWidth {
			right = ansi.Cut(line, cursorCol+1, lineWidth)
		}

		// Keep the cell's syntax colors; only apply reverse video when the cursor is visible.
		cursorCell := cell
		if !m.cursor.virtualCursor.IsBlinked {
			reverseOn := ansi.Style{}.Reverse(true).String()
			reverseOff := ansi.Style{}.Reverse(false).String()
			cursorCell = reverseOn + cell + reverseOff
			// Ensure reverse stays active even if the cell contains reset sequences.
			cursorCell = strings.ReplaceAll(cursorCell, ansi.ResetStyle, ansi.ResetStyle+reverseOn)
			cursorCell = strings.ReplaceAll(cursorCell, "\x1b[0m", "\x1b[0m"+reverseOn)
		}

		lines[m.cursor.row] = left + cursorCell + right
	}

	// Reassemble and render through the viewport.
	result := strings.Join(lines, "\n")
	m.viewport.SetContent(result)
	return m.viewport.View()
}
