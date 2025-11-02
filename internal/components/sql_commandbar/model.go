package sqlcommandbar

import (
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type SQLCommandBarModel struct {
	registry *database.DBRegistry
	textarea textarea.Model
	width    int
	height   int
	err      error
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

// SetSize updates the dimensions of the SQL command bar
func (m *SQLCommandBarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textarea.SetWidth(width - 2)
	m.textarea.SetHeight(height - 4)
}
