package sqlcommandbar

import (
	"fmt"
)

func (m SQLCommandBarModel) View() string {
	return fmt.Sprintf(
		"SQL Query:\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit)",
	)
}
