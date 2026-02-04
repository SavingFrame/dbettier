package editor

type buffer struct {
	lines []string
}

func (b *buffer) handleBackspace(cursor *editorCursor) {
	lineIdx := cursor.row
	line := b.lines[lineIdx]
	newLine := line[:cursor.col-1] + line[cursor.col:]
	b.lines[lineIdx] = newLine
	cursor.moveLeft(1)
}
