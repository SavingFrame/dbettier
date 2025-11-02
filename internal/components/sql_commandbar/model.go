package sqlcommandbar

import (
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type SQLCommandBarModel struct {
	registry *database.DBRegistry
	textarea textarea.Model
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
	}
}

func (m SQLCommandBarModel) Init() tea.Cmd {
	return textarea.Blink
}
