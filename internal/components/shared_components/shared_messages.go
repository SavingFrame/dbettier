package sharedcomponents

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
)

type QueryCompiler interface {
	Compile() string
	HandleSortChange(orderBy []OrderByClause) QueryCompiler
	GetSortOrders() []OrderByClause
	SetSQLResult(*SQLResultMsg) *SQLResult
	GetSQLResult() *SQLResult
	HasPreviousPage() bool
	HasNextPage() bool
	NextPage() tea.Cmd
	PreviousPage() tea.Cmd
	Rows() [][]any
	PageOffset() int
}

type SQLResult struct {
	Rows          [][]any
	Columns       []string // Maybe change, set types for columns, etc
	Total         int      // Total rows available (for pagination)
	TotalFetched  int      // Total rows fetched in this result
	CanFetchTotal bool     // Whether more rows can be fetched
}

type OpenTableMsg struct {
	Table      *database.Table
	DatabaseID string
}

type ExecuteSQLTextMsg struct {
	Query      string
	DatabaseID string
}

type SQLResultMsg struct {
	Rows       [][]any
	Columns    []string // Maybe change, set types for columns, etc
	Query      QueryCompiler
	DatabaseID string
}

type OrderByClause struct {
	ColumnName string
	Direction  string // "ASC" or "DESC"
}

// TableLoadingMsg signals that table data is being loaded
type TableLoadingMsg struct{}

type ComponentTarget int

const (
	TargetSQLCommandBar ComponentTarget = 1 << iota
	TargetTableView
	TargetDBTree
)

var MessageRoutes = map[string]ComponentTarget{
	"sharedcomponents.ExecuteSQLTextMsg":    TargetSQLCommandBar,
	"sharedcomponents.SQLResultMsg":         TargetTableView | TargetSQLCommandBar,
	"sharedcomponents.OrderByChangeMsg":     TargetSQLCommandBar,
	"sharedcomponents.OpenTableMsg":         TargetSQLCommandBar,
	"sharedcomponents.ReapplyTableQueryMsg": TargetSQLCommandBar,
	"sharedcomponents.TableLoadingMsg":      TargetTableView,
}

func GetMessageType(msg tea.Msg) string {
	return fmt.Sprintf("%T", msg)
}
