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

func (b *buffer) handleSpace(cursor *editorCursor) {
	lineIdx := cursor.row
	line := b.lines[lineIdx]
	newLine := line[:cursor.col] + " " + line[cursor.col:]
	b.lines[lineIdx] = newLine
	cursor.moveRight(1, b)
}

func (b *buffer) handleCharacterInput(cursor *editorCursor, char string) {
	line := b.lines[cursor.row]
	line = line[:cursor.col] + char + line[cursor.col:]
	b.lines[cursor.row] = line
	cursor.moveRight(1, b)
}
