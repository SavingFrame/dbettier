package components

import (
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/components/dbtree"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	sqlcommandbar "github.com/SavingFrame/dbettier/internal/components/sql_commandbar"
	"github.com/SavingFrame/dbettier/internal/components/tableview"
	"github.com/SavingFrame/dbettier/internal/database"
)

type FocusedPane int

const (
	FocusDBTree FocusedPane = iota
	FocusTableView
	FocusSQLCommandBar
)

var paneOrder = []FocusedPane{FocusDBTree, FocusTableView, FocusSQLCommandBar}

const (
	DBTreeWidthRatio    = 0.20 // 35% of screen width for dbtree
	SQLCommandBarHeight = 30   // lines
)

type rootScreenModel struct {
	dbtree        dbtree.DBTreeModel
	tableview     tableview.TableViewModel
	sqlCommandBar sqlcommandbar.SQLCommandBarModel

	// State
	focusedPane  FocusedPane
	notification *notifications.Notification
	width        int
	height       int
	registry     *database.DBRegistry
}

func RootScreen(registry *database.DBRegistry) rootScreenModel {
	// Initialize all three components for split layout
	return rootScreenModel{
		dbtree:        dbtree.DBTreeScreen(registry),
		tableview:     tableview.TableViewScreen(),
		sqlCommandBar: sqlcommandbar.SQLCommandBarScreen(registry),
		focusedPane:   FocusDBTree,
		registry:      registry,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	log.Println("RootScreenModel Init() called")
	var cmds []tea.Cmd
	cmds = append(cmds, m.sqlCommandBar.InitialSQLCommand())
	switch m.focusedPane {
	case FocusDBTree:
		cmds = append(cmds, m.dbtree.Init())
	case FocusTableView:
		cmds = append(cmds, m.tableview.Init())
	case FocusSQLCommandBar:
		cmds = append(cmds, m.sqlCommandBar.Init())
	}
	return tea.Batch(cmds...)
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case notifications.ShowNotificationMsg:
		m.notification = &notifications.Notification{
			Message: msg.Message,
			Level:   msg.Level,
		}
		return m, tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
			return notifications.ClearNotificationMsg{}
		})
	case notifications.ClearNotificationMsg:
		m.notification = nil
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate dimensions for each component
		leftWidth := int(float64(m.width) * DBTreeWidthRatio)
		rightWidth := m.width - leftWidth
		topHeight := m.height - SQLCommandBarHeight

		// Update component sizes (accounting for borders: 2 per side = 4 for width, 2 for height)
		// Each border style will add 2 to width and 2 to height, so we subtract those
		m.dbtree.SetSize(leftWidth-4, m.height-4)
		m.tableview.SetSize(rightWidth-4, topHeight-4)
		m.sqlCommandBar.SetSize(rightWidth-4, SQLCommandBarHeight-4)
		return m, nil

	case tea.KeyMsg:
		// Handle focus switching with ctrl+h and ctrl+l
		switch msg.String() {
		case "ctrl+h":
			oldFocus := m.focusedPane
			m.focusedPane = FocusDBTree
			if oldFocus == FocusSQLCommandBar {
				m.sqlCommandBar.Blur()
			}
			return m, nil
		case "ctrl+l":
			oldFocus := m.focusedPane
			oldFocusIndex := 0
			for i, pane := range paneOrder {
				if pane == oldFocus {
					oldFocusIndex = i
					break
				}
			}
			m.focusedPane = paneOrder[(oldFocusIndex+1)%len(paneOrder)]
			if oldFocus != FocusSQLCommandBar {
				return m, m.sqlCommandBar.Focus()
			}
			return m, nil
		case "ctrl+k":
			oldFocus := m.focusedPane
			oldFocusIndex := 0
			for i, pane := range paneOrder {
				if pane == oldFocus {
					oldFocusIndex = i
					break
				}
			}

			m.focusedPane = paneOrder[(oldFocusIndex-1+len(paneOrder))%len(paneOrder)]
			if oldFocus == FocusSQLCommandBar {
				m.sqlCommandBar.Blur()
			}
			return m, nil
		}
	default:
		routedCmds := m.routeToComponents(msg)
		if len(routedCmds) > 0 {
			return m, tea.Batch(routedCmds...)
		}
	}

	// Route to focused pane
	switch m.focusedPane {
	case FocusDBTree:
		var treeModel tea.Model
		treeModel, cmd = m.dbtree.Update(msg)
		m.dbtree = treeModel.(dbtree.DBTreeModel)
		cmds = append(cmds, cmd)
	case FocusTableView:
		var tableModel tea.Model
		tableModel, cmd = m.tableview.Update(msg)
		m.tableview = tableModel.(tableview.TableViewModel)
		cmds = append(cmds, cmd)
	case FocusSQLCommandBar:
		var sqlModel tea.Model
		sqlModel, cmd = m.sqlCommandBar.Update(msg)
		m.sqlCommandBar = sqlModel.(sqlcommandbar.SQLCommandBarModel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *rootScreenModel) routeToComponents(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	msgType := sharedcomponents.GetMessageType(msg)
	targets, shouldRoute := sharedcomponents.MessageRoutes[msgType]
	log.Printf("Routing message of type %s to targets: %d (shouldRoute=%v)\n", msgType, targets, shouldRoute)

	if !shouldRoute {
		return cmds
	}

	if targets&sharedcomponents.TargetDBTree != 0 {
		var treeModel tea.Model
		treeModel, cmd = m.dbtree.Update(msg)
		m.dbtree = treeModel.(dbtree.DBTreeModel)
		cmds = append(cmds, cmd)
	}

	if targets&sharedcomponents.TargetTableView != 0 {
		var tableView tea.Model
		tableView, cmd = m.tableview.Update(msg)
		m.tableview = tableView.(tableview.TableViewModel)
		cmds = append(cmds, cmd)
	}

	if targets&sharedcomponents.TargetSQLCommandBar != 0 {
		var sqlModel tea.Model
		sqlModel, cmd = m.sqlCommandBar.Update(msg)
		m.sqlCommandBar = sqlModel.(sqlcommandbar.SQLCommandBarModel)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m rootScreenModel) View() tea.View {
	var v tea.View
	v.AltScreen = true

	baseView := m.renderSplitLayout()

	if m.notification == nil {
		v.SetContent(baseView)
		return v
	}

	// Create the notification view
	style := m.notification.GetStyle()
	notifView := style.Render(m.notification.Message)

	if m.width > 0 && m.height > 0 {
		// Use Canvas and Layers for compositing notification overlay
		notifWidth := lipgloss.Width(notifView)
		canvas := lipgloss.NewCanvas(
			lipgloss.NewLayer(baseView),
			lipgloss.NewLayer(notifView).X(m.width-notifWidth).Y(0),
		)
		v.SetContent(canvas.Render())
		return v
	}

	v.SetContent(notifView + "\n" + baseView)
	return v
}

func (m rootScreenModel) renderSplitLayout() string {
	if m.width == 0 || m.height == 0 {
		return "Resizing..."
	}

	// Get component views
	treeView := m.renderDBTree()
	tableView := m.renderTableView()
	sqlView := m.renderSQLCommandBar()

	// Compose right column (table view on top, SQL command bar on bottom)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, tableView, sqlView)

	// Compose full layout (tree on left, right column on right)
	layout := lipgloss.JoinHorizontal(lipgloss.Left, treeView, rightColumn)

	return layout
}

func (m rootScreenModel) renderDBTree() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusDBTree {
		borderColor = lipgloss.Color("205")
	}

	// Calculate fixed width for dbtree
	leftWidth := int(float64(m.width) * DBTreeWidthRatio)

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(leftWidth - 4). // Subtract 4 for border padding
		Height(m.height - 4)  // Subtract 4 for border padding

	content := m.dbtree.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderTableView() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusTableView {
		borderColor = lipgloss.Color("205")
	}

	// Don't set explicit height - let the content determine it
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor)

	content := m.tableview.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderSQLCommandBar() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusSQLCommandBar {
		borderColor = lipgloss.Color("205")
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor)

	content := m.sqlCommandBar.RenderContent()
	return borderStyle.Render(content)
}
