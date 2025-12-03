package sqlcommandbar

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	"github.com/SavingFrame/dbettier/internal/database"
)

type SQLCommandBarModel struct {
	registry   *database.DBRegistry
	textarea   textarea.Model
	width      int
	height     int
	err        error
	query      sharedcomponents.SQLQuery
	databaseID string
}

func SQLCommandBarScreen(registry *database.DBRegistry) SQLCommandBarModel {
	ti := textarea.New()
	ti.Placeholder = "Enter SQL command here..."
	ti.ShowLineNumbers = true
	ti.Focus()
	return SQLCommandBarModel{
		registry: registry,
		textarea: ti,
		width:    80,
		height:   30,
	}
}

func (m SQLCommandBarModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m SQLCommandBarModel) InitialSQLCommand() tea.Cmd {
	return func() tea.Msg {
		q := `
SELECT 
    schemaname AS "Schema",
    relname AS "Table Name",
    n_live_tup AS "Live Rows",
    n_dead_tup AS "Dead Rows",
    last_vacuum AS "Last Vacuum",
    last_autovacuum AS "Last Auto Vacuum",
    seq_scan AS "Sequential Scans",
    idx_scan AS "Index Scans"
FROM pg_stat_user_tables
		`
		return sharedcomponents.SetSQLTextMsg{
			Query: sharedcomponents.SQLQuery{
				BaseQuery: q,
				SortOrders: []sharedcomponents.OrderByClause{
					{
						ColumnName: "n_live_tup",
						Direction:  "DESC",
					},
				},
			},
			DatabaseID: m.registry.GetAll()[0].ID,
		}
	}
}

// SetSize updates the dimensions of the SQL command bar
func (m *SQLCommandBarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textarea.SetWidth(width - 2)
	m.textarea.SetHeight(height - 1)
}
