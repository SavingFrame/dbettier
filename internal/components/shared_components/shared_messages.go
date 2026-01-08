package sharedcomponents

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
)

type QueryCompiler interface {
	Compile() string
	HandleSortChange(orderBy OrderByClauses) tea.Cmd
	GetSortOrders() OrderByClauses
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
type OrderByClauses []OrderByClause

func (o OrderByClauses) String() string {
	strs := make([]string, len(o))
	for i, clause := range o {
		strs[i] = fmt.Sprintf("\"%s\" %s", clause.ColumnName, clause.Direction)
	}
	return strings.Join(strs, ", ")
}

func ParseOrderByClauses(s string) (OrderByClauses, error) {
	var clauses OrderByClauses
	for orderClause := range strings.SplitSeq(s, ",") {
		orderClause = strings.TrimSpace(orderClause)
		parts := strings.SplitN(orderClause, " ", 2)
		var columnName, direction string
		switch len(parts) {
		case 2:
			columnName = strings.Trim(parts[0], "\"")
			direction = strings.ToUpper(strings.TrimSpace(parts[1]))
		case 1:
			columnName = strings.Trim(parts[0], "\"")
			direction = "ASC"
		default:
			return nil, fmt.Errorf("invalid order by clause: %s", orderClause)
		}
		clauses = append(clauses, OrderByClause{
			ColumnName: columnName,
			Direction:  direction,
		})
	}
	return clauses, nil
}

type TableLoadingMsg struct{}

// LogLevel defines the severity of a log entry
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogSuccess
	LogWarning
	LogError
	LogSQL
)

// AddLogMsg is a message to add a log entry to the log panel
type AddLogMsg struct {
	Message string
	Level   LogLevel
}

type ComponentTarget int

const (
	TargetSQLCommandBar ComponentTarget = 1 << iota
	TargetTableView
	TargetDBTree
	TargetLogPanel
	TargetWorkspace
)

var MessageRoutes = map[string]ComponentTarget{
	"sharedcomponents.ExecuteSQLTextMsg":    TargetWorkspace,
	"sharedcomponents.SQLResultMsg":         TargetTableView | TargetSQLCommandBar,
	"sharedcomponents.OrderByChangeMsg":     TargetSQLCommandBar,
	"sharedcomponents.OpenTableMsg":         TargetWorkspace,
	"sharedcomponents.ReapplyTableQueryMsg": TargetWorkspace,
	"sharedcomponents.TableLoadingMsg":      TargetTableView,
	"sharedcomponents.AddLogMsg":            TargetLogPanel,
}

func GetMessageType(msg tea.Msg) string {
	return fmt.Sprintf("%T", msg)
}
