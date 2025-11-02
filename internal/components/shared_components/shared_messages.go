package sharedcomponents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type SetSQLTextMsg struct {
	Command    string
	DatabaseID string
}

type ComponentTarget int

const (
	TargetSQLCommandBar ComponentTarget = 1 << iota
	TargetTableView
	TargetDBTree
)

var MessageRoutes = map[string]ComponentTarget{
	"sharedcomponents.SetSQLTextMsg": TargetSQLCommandBar | TargetTableView,
}

func GetMessageType(msg tea.Msg) string {
	return fmt.Sprintf("%T", msg)
}
