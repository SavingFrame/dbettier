package sharedcomponents

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type ComponentTarget int

const (
	TargetSQLCommandBar ComponentTarget = 1 << iota
	TargetTableView
	TargetDBTree
	TargetLogPanel
	TargetWorkspace
	TargetStatusBar
)

var MessageRoutes = map[string]ComponentTarget{
	"messages.ExecuteSQLTextMsg":      TargetWorkspace,
	"query.SQLResultMsg":              TargetTableView | TargetSQLCommandBar,
	"messages.OpenTableAndExecuteMsg": TargetWorkspace,
	"query.ReapplyTableQueryMsg":      TargetWorkspace,
	"messages.TableLoadingMsg":        TargetTableView,
	"messages.AddLogMsg":              TargetLogPanel,
	"editor.EditorModeChangedMsg":     TargetStatusBar,
	"editor.EditorCursorMovedMsg":     TargetStatusBar,
	"messages.OpenQueryTabMsg":        TargetWorkspace,
	"query.UpdateTableMsg":            TargetTableView,
}

func GetMessageType(msg tea.Msg) string {
	return fmt.Sprintf("%T", msg)
}
