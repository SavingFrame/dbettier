// Package workspace provides the workspace component that manages multiple tabs,
// each containing a tableview and SQL command bar.
package workspace

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	sqlcommandbarv2 "github.com/SavingFrame/dbettier/internal/components/sql_commandbar_v2"
	"github.com/SavingFrame/dbettier/internal/components/tableview"
	"github.com/SavingFrame/dbettier/internal/database"
)

// TabType represents the type of content in a tab
type TabType int

const (
	TabTypeQuery TabType = iota
	TabTypeTable
)

type TabSize struct {
	width  int
	height int
}

// Tab represents a single tab with its own tableview and sqlcommandbar state
type Tab struct {
	ID            string
	Name          string
	Type          TabType
	TableView     tableview.TableViewModel
	SQLCommandBar sqlcommandbarv2.SQLCommandBarModel
	DatabaseID    string
}

// Icon returns the nerd font icon for the tab type
func (t Tab) Icon() string {
	switch t.Type {
	case TabTypeTable:
		return "󰓫"
	case TabTypeQuery:
		return "󰆍"
	default:
		return "󰆍"
	}
}

// Workspace manages the workspace with multiple tabs
type Workspace struct {
	tabs         []Tab
	activeIndex  int
	width        int
	height       int
	queryCounter int
	registry     *database.DBRegistry

	// Scroll state for tab overflow
	scrollOffset int

	SQLCommandBarSize TabSize
	TableViewSize     TabSize
}

// New creates a new workspace with an initial query tab
func New(registry *database.DBRegistry) Workspace {
	m := Workspace{
		tabs:         []Tab{},
		activeIndex:  0,
		queryCounter: 0,
		registry:     registry,
		scrollOffset: 0,
	}
	// Create initial query tab
	// m.addQueryTab()
	return m
}

// AddQueryTab creates a new query tab
func (w *Workspace) AddQueryTab(databaseID string) {
	w.queryCounter++
	tab := Tab{
		ID:            fmt.Sprintf("query-%d", w.queryCounter),
		Name:          fmt.Sprintf("Query %d", w.queryCounter),
		Type:          TabTypeQuery,
		TableView:     tableview.TableViewScreen(),
		SQLCommandBar: sqlcommandbarv2.NewSQLCommandBarModel(nil, w.registry, databaseID, false), // TODO: Fix this shit
	}
	tab.TableView.SetSize(w.TableViewSize.width, w.TableViewSize.height)
	tab.SQLCommandBar.SetSize(w.SQLCommandBarSize.width, w.SQLCommandBarSize.height)
	w.tabs = append(w.tabs, tab)
	w.activeIndex = len(w.tabs) - 1
	w.ensureActiveTabVisible()
}

// AddTableTab creates a new tab for a table
func (w *Workspace) AddTableTab(tableName string, databaseID string) int {
	tab := Tab{
		ID:            fmt.Sprintf("table-%s-%s", databaseID, tableName),
		Name:          tableName,
		Type:          TabTypeTable,
		TableView:     tableview.TableViewScreen(),
		SQLCommandBar: sqlcommandbarv2.NewSQLCommandBarModel(nil, w.registry, databaseID, true),
	}
	tab.TableView.SetSize(w.TableViewSize.width, w.TableViewSize.height)
	tab.SQLCommandBar.SetSize(w.SQLCommandBarSize.width, w.SQLCommandBarSize.height)
	w.tabs = append(w.tabs, tab)
	w.activeIndex = len(w.tabs) - 1
	w.ensureActiveTabVisible()
	return w.activeIndex
}

// Tabs returns all tabs
func (w *Workspace) Tabs() []Tab {
	return w.tabs
}

// ActiveIndex returns the currently active tab index
func (w *Workspace) ActiveIndex() int {
	return w.activeIndex
}

// ActiveTab returns the currently active tab
func (w *Workspace) ActiveTab() *Tab {
	if w.activeIndex >= 0 && w.activeIndex < len(w.tabs) {
		return &w.tabs[w.activeIndex]
	}
	return nil
}

// SetActiveIndex sets the active tab by index
func (w *Workspace) SetActiveIndex(index int) {
	if index >= 0 && index < len(w.tabs) {
		w.activeIndex = index
		w.ensureActiveTabVisible()
	}
}

// NextTab switches to the next tab
func (w *Workspace) NextTab() {
	if len(w.tabs) > 0 {
		w.activeIndex = (w.activeIndex + 1) % len(w.tabs)
		w.ensureActiveTabVisible()
	}
}

// PrevTab switches to the previous tab
func (w *Workspace) PrevTab() {
	if len(w.tabs) > 0 {
		w.activeIndex = (w.activeIndex - 1 + len(w.tabs)) % len(w.tabs)
		w.ensureActiveTabVisible()
	}
}

// CloseTab closes a tab by index
func (w *Workspace) CloseTab(index int) {
	if index < 0 || index >= len(w.tabs) {
		return
	}

	// Remove the tab
	w.tabs = append(w.tabs[:index], w.tabs[index+1:]...)

	// Adjust active index if needed
	if len(w.tabs) == 0 {
		// Create a new query tab if all tabs are closed
		// w.addQueryTab()
		return
	}

	if w.activeIndex >= len(w.tabs) {
		w.activeIndex = len(w.tabs) - 1
	} else if w.activeIndex > index {
		w.activeIndex--
	}
	w.ensureActiveTabVisible()
}

// CloseActiveTab closes the currently active tab
func (w *Workspace) CloseActiveTab() {
	w.CloseTab(w.activeIndex)
}

// SetSize updates the dimensions of the tab bar
func (w *Workspace) SetSize(width, height int) {
	w.width = width
	w.height = height
}

// SetTabSizes updates the sizes of tableview and sqlcommandbar for all tabs
func (w *Workspace) SetTabSizes(tableWidth, tableHeight, sqlWidth, sqlHeight int) {
	w.SQLCommandBarSize = TabSize{width: sqlWidth, height: sqlHeight}
	w.TableViewSize = TabSize{width: tableWidth, height: tableHeight}
	for i := range w.tabs {
		w.tabs[i].TableView.SetSize(tableWidth, tableHeight)
		w.tabs[i].SQLCommandBar.SetSize(sqlWidth, sqlHeight)
	}
}

// Init initializes the active tab's components
func (w Workspace) Init() tea.Cmd {
	if tab := w.ActiveTab(); tab != nil {
		return tea.Batch(
			tab.TableView.Init(),
			tab.SQLCommandBar.Init(),
		)
	}
	return nil
}

// InitialSQLCommand returns the initial SQL command for the first tab
func (w Workspace) InitialSQLCommand() tea.Cmd {
	return nil
	// TODO: Temperary disabled, we should move this logic from SQLCommandBar to Workspace

	// if tab := w.ActiveTab(); tab != nil {
	// 	return tab.SQLCommandBar.InitialSQLCommand()
	// }
	// return nil
}

// UpdateActiveTableView updates the active tab's tableview
func (w *Workspace) UpdateActiveTableView(msg tea.Msg) tea.Cmd {
	if tab := w.ActiveTab(); tab != nil {
		log.Printf("Routing message to active tab's TableView: %+v", msg)
		model, cmd := tab.TableView.Update(msg)
		tab.TableView = model.(tableview.TableViewModel)
		return cmd
	}
	return nil
}

// UpdateActiveSQLCommandBar updates the active tab's sqlcommandbar
func (w *Workspace) UpdateActiveSQLCommandBar(msg tea.Msg) tea.Cmd {
	if tab := w.ActiveTab(); tab != nil {
		model, cmd := tab.SQLCommandBar.Update(msg)
		tab.SQLCommandBar = model.(sqlcommandbarv2.SQLCommandBarModel)
		return cmd
	}
	return nil
}

// Focus focuses the active tab's sqlcommandbar
func (w *Workspace) Focus() tea.Cmd {
	if tab := w.ActiveTab(); tab != nil {
		return tab.SQLCommandBar.Focus()
	}
	return nil
}

// Blur blurs the active tab's sqlcommandbar
func (w *Workspace) Blur() {
	if tab := w.ActiveTab(); tab != nil {
		tab.SQLCommandBar.Blur()
	}
}

// GetActiveTableViewSize returns the size of the active tableview
func (w *Workspace) GetActiveTableViewSize() (int, int) {
	if tab := w.ActiveTab(); tab != nil {
		return tab.TableView.GetSize()
	}
	return 0, 0
}

// RenderActiveTableView returns the rendered content of the active tableview
func (w *Workspace) RenderActiveTableView() string {
	if tab := w.ActiveTab(); tab != nil {
		return tab.TableView.RenderContent()
	}
	return ""
}

// RenderActiveSQLCommandBar returns the rendered content of the active sqlcommandbar
func (w *Workspace) RenderActiveSQLCommandBar() string {
	if tab := w.ActiveTab(); tab != nil {
		return tab.SQLCommandBar.RenderContent()
	}
	return ""
}

// Width returns the tab bar width
func (w *Workspace) Width() int {
	return w.width
}

// ScrollOffset returns the current scroll offset
func (w *Workspace) ScrollOffset() int {
	return w.scrollOffset
}

// ensureActiveTabVisible adjusts scroll offset to keep active tab visible
func (w *Workspace) ensureActiveTabVisible() {
	// This will be properly calculated in the view based on actual tab widths
	// For now, simple logic: if active tab is before scroll offset, scroll to it
	if w.activeIndex < w.scrollOffset {
		w.scrollOffset = w.activeIndex
	}
}

// SetScrollOffset sets the scroll offset (used by view calculations)
func (w *Workspace) SetScrollOffset(offset int) {
	w.scrollOffset = offset
}

// HandleOpenTable handles opening a table in a new tab
func (w *Workspace) HandleOpenTable(msg sharedcomponents.OpenTableMsg) tea.Cmd {
	// Create new table tab
	w.AddTableTab(msg.Table.Name, msg.DatabaseID)

	// Return command to execute the table query on the new tab
	return func() tea.Msg {
		return sharedcomponents.OpenTableMsg{
			Table:      msg.Table,
			DatabaseID: msg.DatabaseID,
		}
	}
}
