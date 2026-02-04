package editor

import (
	"bytes"

	"github.com/SavingFrame/dbettier/internal/theme"
	"github.com/alecthomas/chroma/v2/quick"
)

func highlightCode(content string) string {
	buf := new(bytes.Buffer)
	s := theme.Current().Name
	if err := quick.Highlight(buf, content, "sql", "terminal256", s); err != nil {
		return content
	}
	return buf.String()
}
