package components

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/components/dbtree"
	"github.com/SavingFrame/dbettier/internal/components/logpanel"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"
	sqlcommandbar "github.com/SavingFrame/dbettier/internal/components/sql_commandbar"
	"github.com/SavingFrame/dbettier/internal/components/tableview"
	"github.com/SavingFrame/dbettier/internal/database"
	zone "github.com/lrstanley/bubblezone/v2"
)

type FocusedPane int

const (
	FocusDBTree FocusedPane = iota
	FocusTableView
	FocusSQLCommandBar
	FocusLogPanel
)

var paneOrder = []FocusedPane{FocusDBTree, FocusTableView, FocusSQLCommandBar, FocusLogPanel}

const (
	DBTreeWidthRatio         = 0.20 // 20% of screen width for dbtree
	SQLCommandBarHeightRatio = 30   // percent
	SQLCommandBarWidthRatio  = 0.60 // 60% of bottom row for SQL command bar
)

type rootScreenModel struct {
	dbtree        dbtree.DBTreeModel
	tableview     tableview.TableViewModel
	sqlCommandBar sqlcommandbar.SQLCommandBarModel
	logPanel      logpanel.LogPanelModel

	// State
	focusedPane  FocusedPane
	notification *notifications.Notification
	width        int
	height       int
	registry     *database.DBRegistry

	// Help
	help help.Model
	keys GlobalKeyMap
}

func RootScreen(registry *database.DBRegistry) rootScreenModel {
	// Initialize all four components for split layout
	return rootScreenModel{
		dbtree:        dbtree.DBTreeScreen(registry),
		tableview:     tableview.TableViewScreen(),
		sqlCommandBar: sqlcommandbar.SQLCommandBarScreen(registry),
		logPanel:      logpanel.LogPanelScreen(),
		focusedPane:   FocusDBTree,
		registry:      registry,
		help:          help.New(),
		keys:          DefaultGlobalKeyMap,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.sqlCommandBar.InitialSQLCommand())
	switch m.focusedPane {
	case FocusDBTree:
		cmds = append(cmds, m.dbtree.Init())
	case FocusTableView:
		cmds = append(cmds, m.tableview.Init())
	case FocusSQLCommandBar:
		cmds = append(cmds, m.sqlCommandBar.Init())
	case FocusLogPanel:
		cmds = append(cmds, m.logPanel.Init())
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

	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return m, nil
		}
		if zone.Get("dbTree").InBounds(msg) {
			m.focusedPane = FocusDBTree
		} else if zone.Get("tableview").InBounds(msg) {
			m.focusedPane = FocusTableView
		} else if zone.Get("sqlCommandBar").InBounds(msg) {
			m.focusedPane = FocusSQLCommandBar
		} else if zone.Get("logPanel").InBounds(msg) {
			m.focusedPane = FocusLogPanel
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update help width for proper truncation
		m.help.SetWidth(msg.Width)

		// Reserve space for short help bar at bottom (1 line)
		helpHeight := 1
		availableHeight := m.height - helpHeight

		// Calculate dimensions for each component
		leftWidth := int(float64(m.width) * DBTreeWidthRatio)
		rightWidth := m.width - leftWidth
		bottomRowHeight := int(float64(availableHeight) * (float64(SQLCommandBarHeightRatio) / 100.0))
		tableViewHeight := availableHeight - bottomRowHeight

		// Calculate widths for SQL command bar and log panel (split bottom row)
		sqlCommandBarWidth := int(float64(rightWidth) * SQLCommandBarWidthRatio)
		logPanelWidth := rightWidth - sqlCommandBarWidth

		// Update component sizes (accounting for borders: 2 per side = 4 for width, 2 for height)
		// Each border style will add 2 to width and 2 to height, so we subtract those
		m.dbtree.SetSize(leftWidth-4, availableHeight-2)
		m.tableview.SetSize(rightWidth-4, tableViewHeight-2)
		m.sqlCommandBar.SetSize(sqlCommandBarWidth-4, bottomRowHeight)
		m.logPanel.SetSize(logPanelWidth+4, bottomRowHeight-2) // TODO: Something weird with width here. if i add +4 it show more content?
		return m, nil

	case tea.KeyMsg:
		// Handle help toggle first
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}

		// Handle focus switching with ctrl+h and ctrl+l
		switch {
		case key.Matches(msg, m.keys.FocusLeft):
			oldFocus := m.focusedPane
			m.focusedPane = FocusDBTree
			if oldFocus == FocusSQLCommandBar {
				m.sqlCommandBar.Blur()
			}
			return m, nil
		case key.Matches(msg, m.keys.FocusNext):
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
		case key.Matches(msg, m.keys.FocusPrev):
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
		case key.Matches(msg, m.keys.Escape) && m.help.ShowAll:
			m.help.ShowAll = false
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
	case FocusLogPanel:
		var logModel tea.Model
		logModel, cmd = m.logPanel.Update(msg)
		m.logPanel = logModel.(logpanel.LogPanelModel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *rootScreenModel) routeToComponents(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	msgType := sharedcomponents.GetMessageType(msg)
	targets, shouldRoute := sharedcomponents.MessageRoutes[msgType]

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

	if targets&sharedcomponents.TargetLogPanel != 0 {
		var logModel tea.Model
		logModel, cmd = m.logPanel.Update(msg)
		m.logPanel = logModel.(logpanel.LogPanelModel)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m rootScreenModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	baseView := m.renderSplitLayout()

	// Render short help bar at the bottom (always visible)
	shortHelpView := m.help.ShortHelpView(m.getContextKeyMap().ShortHelp())
	fullView := lipgloss.JoinVertical(lipgloss.Left, baseView, shortHelpView)

	// If full help is toggled, render it as a centered popup overlay
	if m.help.ShowAll && m.width > 0 && m.height > 0 {
		fullView = m.renderWithHelpPopup(fullView)
	}

	if m.notification == nil {
		v.SetContent(zone.Scan(fullView))
		return v
	}

	// Create the notification view
	style := m.notification.GetStyle()
	notifView := style.Render(m.notification.Message)

	if m.width > 0 && m.height > 0 {
		// Use Canvas and Layers for compositing notification overlay
		notifWidth := lipgloss.Width(notifView)
		canvas := lipgloss.NewCanvas(
			lipgloss.NewLayer(fullView),
			lipgloss.NewLayer(notifView).X(m.width-notifWidth).Y(0),
		)
		v.SetContent(zone.Scan(canvas.Render()))
		return v
	}

	v.SetContent(zone.Scan(notifView + "\n" + fullView))
	return v
}

// renderWithHelpPopup renders the full help as a centered popup overlay
func (m rootScreenModel) renderWithHelpPopup(baseView string) string {
	fullHelpContent := m.help.FullHelpView(m.getContextKeyMap().FullHelp())

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)

	title := titleStyle.Render("Keyboard Shortcuts")
	helpPopup := popupStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, fullHelpContent))

	popupWidth := lipgloss.Width(helpPopup)
	popupHeight := lipgloss.Height(helpPopup)
	x := (m.width - popupWidth) / 2
	y := (m.height - popupHeight) / 2
	x = max(0, x)
	y = max(0, y)

	canvas := lipgloss.NewCanvas(
		lipgloss.NewLayer(baseView),
		lipgloss.NewLayer(helpPopup).X(x).Y(y),
	)

	return canvas.Render()
}

func (m rootScreenModel) renderSplitLayout() string {
	if m.width == 0 || m.height == 0 {
		return "Resizing..."
	}

	// Get component views
	treeView := zone.Mark("dbTree", m.renderDBTree())
	tableView := zone.Mark("tableview", m.renderTableView())
	sqlView := zone.Mark("sqlCommandBar", m.renderSQLCommandBar())
	logView := zone.Mark("logPanel", m.renderLogPanel())

	// Compose bottom row (SQL command bar on left, log panel on right)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Left, sqlView, logView)

	// Compose right column (table view on top, bottom row on bottom)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, tableView, bottomRow)

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

	helpHeight := 1

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(leftWidth - 4). // Subtract 4 for border padding
		Height(m.height - helpHeight)

	content := m.dbtree.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderTableView() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusTableView {
		borderColor = lipgloss.Color("205")
	}

	leftWidth := int(float64(m.width) * DBTreeWidthRatio)
	// Don't set explicit height - let the content determine it
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).Width(m.width - leftWidth)

	content := m.tableview.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderSQLCommandBar() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusSQLCommandBar {
		borderColor = lipgloss.Color("205")
	}

	leftWidth := int(float64(m.width) * DBTreeWidthRatio)
	rightWidth := m.width - leftWidth
	sqlCommandBarWidth := int(float64(rightWidth) * SQLCommandBarWidthRatio)
	_, tableViewY := m.tableview.GetSize()

	helpHeight := 1

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(sqlCommandBarWidth - 4).
		Height(m.height - tableViewY - helpHeight)

	content := m.sqlCommandBar.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderLogPanel() string {
	borderColor := lipgloss.Color("240")
	if m.focusedPane == FocusLogPanel {
		borderColor = lipgloss.Color("205")
	}

	leftWidth := int(float64(m.width) * DBTreeWidthRatio)
	rightWidth := m.width - leftWidth
	sqlCommandBarWidth := int(float64(rightWidth) * SQLCommandBarWidthRatio)
	logPanelWidth := rightWidth - sqlCommandBarWidth
	_, tableViewY := m.tableview.GetSize()

	helpHeight := 1

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(logPanelWidth + 4).
		Height(m.height - tableViewY - helpHeight)

	content := m.logPanel.RenderContent()
	return borderStyle.Render(content)
}

// combinedKeyMap combines global keys with focused pane keys for help display
type combinedKeyMap struct {
	global       GlobalKeyMap
	paneKeys     []key.Binding
	fullPaneKeys [][]key.Binding
}

// ShortHelp returns keybindings for the short help view
func (k combinedKeyMap) ShortHelp() []key.Binding {
	bindings := k.paneKeys
	bindings = append(bindings, k.global.Help, k.global.Quit)
	return bindings
}

// FullHelp returns keybindings for the expanded help view
func (k combinedKeyMap) FullHelp() [][]key.Binding {
	result := k.fullPaneKeys
	// Add global keys as the last column
	result = append(result, []key.Binding{k.global.FocusLeft, k.global.FocusNext, k.global.FocusPrev, k.global.Help, k.global.Quit})
	return result
}

// getContextKeyMap returns a combined keymap for the currently focused pane
func (m rootScreenModel) getContextKeyMap() combinedKeyMap {
	combined := combinedKeyMap{
		global: m.keys,
	}

	switch m.focusedPane {
	case FocusDBTree:
		keys := dbtree.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.fullPaneKeys = keys.FullHelp()
	case FocusTableView:
		keys := tableview.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.fullPaneKeys = keys.FullHelp()
	case FocusSQLCommandBar:
		keys := sqlcommandbar.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.fullPaneKeys = keys.FullHelp()
	case FocusLogPanel:
		keys := logpanel.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.fullPaneKeys = keys.FullHelp()
	}

	return combined
}
