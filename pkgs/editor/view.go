package editor

import "strings"

func (m SQLEditor) View() string {
	const placeholder = "\x00CURSOR\x00"

	m.cursor.setCharUnderCursor(m.buffer.lines)
	var s strings.Builder
	for lineNumber, line := range m.buffer.lines {
		if lineNumber == m.cursor.row {
			s.WriteString(line[:m.cursor.col])
			s.WriteString(placeholder)
			if m.cursor.col < len(line) {
				s.WriteString(line[m.cursor.col+1:])
			}
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}

	highlighted := highlightCode(s.String())
	result := strings.Replace(highlighted, placeholder, m.cursor.virtualCursor.View(), 1)

	m.viewport.SetContent(result)
	return m.viewport.View()
}
