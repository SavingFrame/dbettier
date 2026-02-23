// Package components provides the main UI components for the dbettier application.
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
	"github.com/SavingFrame/dbettier/internal/components/statusbar"
	"github.com/SavingFrame/dbettier/internal/components/workspace"
	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/SavingFrame/dbettier/internal/theme"
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

type rootLayout struct {
	helpHeight   int
	statusHeight int

	contentHeight int

	leftWidth  int
	rightWidth int

	tabBarHeight int
	tableHeight  int
	bottomHeight int

	sqlWidth int
	logWidth int
}

type rootScreenModel struct {
	dbtree    dbtree.DBTreeModel
	workspace workspace.Workspace
	logPanel  logpanel.LogPanelModel
	statusBar statusbar.StatusBarModel

	// State
	focusedPane  FocusedPane
	notification *notifications.Notification
	width        int
	height       int
	layout       rootLayout
	registry     *database.DBRegistry

	// Help
	help help.Model
	keys GlobalKeyMap
}

func RootScreen(registry *database.DBRegistry) rootScreenModel {
	// Initialize all components for split layout
	return rootScreenModel{
		dbtree:      dbtree.DBTreeScreen(registry),
		statusBar:   statusbar.NewStatusBarModel(),
		workspace:   workspace.New(registry),
		logPanel:    logpanel.LogPanelScreen(),
		focusedPane: FocusDBTree,
		registry:    registry,
		help:        help.New(),
		keys:        DefaultGlobalKeyMap,
	}
}

func splitByRatio(total int, ratio float64) (first, second int) {
	if total <= 0 {
		return 0, 0
	}
	if total == 1 {
		return 1, 0
	}
	first = int(float64(total) * ratio)
	first = max(1, min(total-1, first))
	second = total - first
	return first, second
}

func borderedInner(total int) int {
	return max(0, total-2)
}

func calculateRootLayout(width, height int) rootLayout {
	layout := rootLayout{
		helpHeight:   1,
		statusHeight: 1,
	}

	if width <= 0 || height <= 0 {
		return layout
	}

	layout.contentHeight = max(0, height-layout.helpHeight-layout.statusHeight)
	layout.tabBarHeight = min(workspace.TabBarHeight, layout.contentHeight)

	layout.leftWidth, layout.rightWidth = splitByRatio(width, DBTreeWidthRatio)

	bodyHeight := max(0, layout.contentHeight-layout.tabBarHeight)
	tableRatio := 1.0 - (float64(SQLCommandBarHeightRatio) / 100.0)
	layout.tableHeight, layout.bottomHeight = splitByRatio(bodyHeight, tableRatio)
	if bodyHeight > 0 && layout.bottomHeight == 0 {
		layout.bottomHeight = 1
		layout.tableHeight = max(0, bodyHeight-layout.bottomHeight)
	}

	layout.sqlWidth, layout.logWidth = splitByRatio(layout.rightWidth, SQLCommandBarWidthRatio)

	return layout
}

func (m rootScreenModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.workspace.InitialSQLCommand())
	switch m.focusedPane {
	case FocusDBTree:
		cmds = append(cmds, m.dbtree.Init())
	case FocusTableView:
		cmds = append(cmds, m.workspace.Init())
	case FocusSQLCommandBar:
		cmds = append(cmds, m.workspace.Init())
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

		if zone.Get("tabbar").InBounds(msg) {
			handled, tabCmd := m.workspace.HandleMouseClick(msg)
			if handled {
				return m, tabCmd
			}
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
		m.layout = calculateRootLayout(msg.Width, msg.Height)

		// Update help width for proper truncation
		m.help.SetWidth(msg.Width)

		m.dbtree.SetSize(
			borderedInner(m.layout.leftWidth),
			borderedInner(m.layout.contentHeight),
		)

		m.workspace.SetSize(m.layout.rightWidth, m.layout.tabBarHeight)
		m.workspace.SetTabSizes(
			borderedInner(m.layout.rightWidth),
			borderedInner(m.layout.tableHeight),
			borderedInner(m.layout.sqlWidth),
			borderedInner(m.layout.bottomHeight),
		)

		m.logPanel.SetSize(
			borderedInner(m.layout.logWidth),
			borderedInner(m.layout.bottomHeight),
		)

		m.statusBar.SetSize(m.width, m.layout.statusHeight)
		return m, nil

	case tea.KeyMsg:
		// Handle help toggle first
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}

		// Handle tab navigation (works globally, not just when focused on tabbar)
		if m.workspace.HandleKeys(msg) {
			return m, nil
		}

		// Handle focus switching with ctrl+h and ctrl+l
		switch {
		case key.Matches(msg, m.keys.FocusLeft):
			oldFocus := m.focusedPane
			m.focusedPane = FocusDBTree
			if oldFocus == FocusSQLCommandBar {
				m.workspace.Blur()
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
			if m.focusedPane == FocusSQLCommandBar {
				return m, m.workspace.Focus()
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
				m.workspace.Blur()
			}
			if m.focusedPane == FocusSQLCommandBar {
				return m, m.workspace.Focus()
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
		cmd = m.workspace.UpdateActiveTableView(msg)
		cmds = append(cmds, cmd)
	case FocusSQLCommandBar:
		cmd = m.workspace.UpdateActiveSQLCommandBar(msg)
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

	if targets&sharedcomponents.TargetLogPanel != 0 {
		var logModel tea.Model
		logModel, cmd = m.logPanel.Update(msg)
		m.logPanel = logModel.(logpanel.LogPanelModel)
		cmds = append(cmds, cmd)
	}

	// Route to active tab's components
	if targets&sharedcomponents.TargetWorkspace != 0 {
		// Handle OpenTableMsg specially - create new tab
		var workspaceModel tea.Model
		workspaceModel, cmd = m.workspace.Update(msg)
		m.workspace = workspaceModel.(workspace.Workspace)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if targets&sharedcomponents.TargetTableView != 0 {
		cmd = m.workspace.UpdateActiveTableView(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if targets&sharedcomponents.TargetSQLCommandBar != 0 {
		cmd = m.workspace.UpdateActiveSQLCommandBar(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if targets&sharedcomponents.TargetStatusBar != 0 {
		var statusBarModel tea.Model
		statusBarModel, cmd = m.statusBar.Update(msg)
		m.statusBar = statusBarModel.(statusbar.StatusBarModel)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m rootScreenModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	colors := theme.Current().Colors
	v.BackgroundColor = colors.Base
	baseHeight := m.layout.contentHeight
	if baseHeight <= 0 {
		baseHeight = max(0, m.height-2)
	}
	baseView := m.renderSplitLayout()
	if m.width > 0 {
		baseView = lipgloss.NewStyle().
			Width(m.width).
			Height(baseHeight).
			Background(colors.Base).
			Render(baseView)
	}

	// Render status bar and short help bar at the bottom (always visible)
	statusBarView := m.statusBar.RenderContent()
	if m.width > 0 {
		statusBarView = lipgloss.NewStyle().
			Width(m.width).
			Height(max(1, m.layout.statusHeight)).
			Background(colors.Surface).
			Render(statusBarView)
	}

	shortHelpView := m.help.ShortHelpView(m.getContextKeyMap().ShortHelp())
	if m.width > 0 {
		shortHelpView = lipgloss.NewStyle().
			Width(m.width).
			Height(max(1, m.layout.helpHeight)).
			Background(colors.Base).
			Render(shortHelpView)
	}

	fullView := lipgloss.JoinVertical(lipgloss.Left, baseView, shortHelpView, statusBarView)

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
	colors := theme.Current().Colors
	fullHelpContent = lipgloss.NewStyle().Background(colors.Surface).Render(fullHelpContent)

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.BorderFocused).
		Background(colors.Surface).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Foreground(colors.Primary).
		Background(colors.Surface).
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
	tabBarView := zone.Mark("tabbar", m.renderTabBar())
	tableView := zone.Mark("tableview", m.renderTableView())
	sqlView := zone.Mark("sqlCommandBar", m.renderSQLCommandBar())
	logView := zone.Mark("logPanel", m.renderLogPanel())

	// Compose bottom row (SQL command bar on left, log panel on right)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Left, sqlView, logView)

	// Compose right column (tab bar on top, table view below, bottom row at bottom)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, tabBarView, tableView, bottomRow)

	// Compose full layout (tree on left, right column on right)
	layout := lipgloss.JoinHorizontal(lipgloss.Left, treeView, rightColumn)

	return layout
}

func (m rootScreenModel) renderDBTree() string {
	colors := theme.Current().Colors
	borderColor := colors.Border
	if m.focusedPane == FocusDBTree {
		borderColor = colors.BorderFocused
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Background(colors.Base).
		Width(m.layout.leftWidth).
		Height(m.layout.contentHeight)

	content := m.dbtree.RenderContent()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderTabBar() string {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Width(m.layout.rightWidth).
		Height(m.layout.tabBarHeight).
		Background(colors.Base).
		Render(m.workspace.RenderTabBar())
}

func (m rootScreenModel) renderTableView() string {
	colors := theme.Current().Colors
	borderColor := colors.Border
	if m.focusedPane == FocusTableView {
		borderColor = colors.BorderFocused
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Background(colors.Base).
		Width(m.layout.rightWidth).
		Height(m.layout.tableHeight)

	content := m.workspace.RenderActiveTableView()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderSQLCommandBar() string {
	colors := theme.Current().Colors
	borderColor := colors.Border
	if m.focusedPane == FocusSQLCommandBar {
		borderColor = colors.BorderFocused
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Background(colors.Base).
		Width(m.layout.sqlWidth).
		Height(m.layout.bottomHeight)

	content := m.workspace.RenderActiveSQLCommandBar()
	return borderStyle.Render(content)
}

func (m rootScreenModel) renderLogPanel() string {
	colors := theme.Current().Colors
	borderColor := colors.Border
	if m.focusedPane == FocusLogPanel {
		borderColor = colors.BorderFocused
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Background(colors.Base).
		Width(m.layout.logWidth).
		Height(m.layout.bottomHeight)

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

	// Tab keys are always available
	tabKeys := workspace.DefaultKeyMap

	switch m.focusedPane {
	case FocusDBTree:
		keys := dbtree.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.paneKeys = append(combined.paneKeys, tabKeys.ShortHelp()...)
		combined.fullPaneKeys = keys.FullHelp()
		combined.fullPaneKeys = append(combined.fullPaneKeys, tabKeys.FullHelp()...)
	case FocusTableView:
		// Combine tableview keys with tab navigation keys
		combined.paneKeys = tabKeys.ShortHelp()
		combined.fullPaneKeys = tabKeys.FullHelp()
	case FocusSQLCommandBar:
		// Combine sqlcommandbar keys with tab navigation keys
		combined.paneKeys = tabKeys.ShortHelp()
		combined.fullPaneKeys = tabKeys.FullHelp()
	case FocusLogPanel:
		keys := logpanel.DefaultKeyMap
		combined.paneKeys = keys.ShortHelp()
		combined.paneKeys = append(combined.paneKeys, tabKeys.ShortHelp()...)
		combined.fullPaneKeys = keys.FullHelp()
		combined.fullPaneKeys = append(combined.fullPaneKeys, tabKeys.FullHelp()...)
	}

	return combined
}
