package components

import (
	"log"
	"time"

	"github.com/SavingFrame/dbettier/internal/components/dbtree"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	sqlcommandbar "github.com/SavingFrame/dbettier/internal/components/sql_commandbar"
	"github.com/SavingFrame/dbettier/internal/components/tableview"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FocusedPane int

const (
	FocusDBTree FocusedPane = iota
	FocusTableView
	FocusSQLCommandBar
)

const (
	DBTreeWidthRatio    = 0.20 // 35% of screen width for dbtree
	SQLCommandBarHeight = 30   // lines
)

type rootScreenModel struct {
	dbtree        dbtree.DBTreeModel
	tableview     tableview.TableViewModel
	sqlCommandBar sqlcommandbar.SQLCommandBarModel

	// Legacy model for backward compatibility (screen switching)
	model tea.Model

	// State
	focusedPane  FocusedPane
	notification *notifications.Notification
	width        int
	height       int
	registry     *database.DBRegistry

	// Layout mode: true = new split layout, false = legacy single model
	useSplitLayout bool
}

func RootScreen(registry *database.DBRegistry) rootScreenModel {
	var rootModel tea.Model

	// Determine if we should use the new split layout
	useSplitLayout := true

	if useSplitLayout {
		// Initialize all three components for split layout
		return rootScreenModel{
			dbtree:         dbtree.DBTreeScreen(registry),
			tableview:      tableview.TableViewScreen(),
			sqlCommandBar:  sqlcommandbar.SQLCommandBarScreen(registry),
			focusedPane:    FocusDBTree,
			registry:       registry,
			useSplitLayout: true,
		}
	} else {
		// Legacy: use DB creator screen
		// screenOne := dbtree.DBTreeScreen(registry)
		screenOne := dbtree.DBTreeScreen(registry)
		rootModel = &screenOne
		return rootScreenModel{
			model:          rootModel,
			registry:       registry,
			useSplitLayout: false,
		}
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	if m.model != nil {
		return m.model.Init()
	}
	return nil
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

		if m.useSplitLayout {
			// Calculate dimensions for each component
			leftWidth := int(float64(m.width) * DBTreeWidthRatio)
			rightWidth := m.width - leftWidth
			topHeight := m.height - SQLCommandBarHeight

			// Update component sizes (accounting for borders: 2 per side = 4 for width, 2 for height)
			// Each border style will add 2 to width and 2 to height, so we subtract those
			m.dbtree.SetSize(leftWidth-4, m.height-4)
			m.tableview.SetSize(rightWidth-4, topHeight-4)
			m.sqlCommandBar.SetSize(rightWidth-4, SQLCommandBarHeight-4)
		}
		return m, nil

	case tea.KeyMsg:
		// Handle focus switching with ctrl+h and ctrl+l
		log.Println("Key pressed:", msg.String())
		if m.useSplitLayout {
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
				m.focusedPane = FocusSQLCommandBar
				if oldFocus != FocusSQLCommandBar {
					return m, m.sqlCommandBar.Focus()
				}
				return m, nil
			case "ctrl+k":
				oldFocus := m.focusedPane
				m.focusedPane = FocusTableView
				if oldFocus == FocusSQLCommandBar {
					m.sqlCommandBar.Blur()
				}
				return m, nil
			}
		}
	}

	// Route updates to appropriate component based on layout mode
	if m.useSplitLayout {
		// Route to focused pane
		switch m.focusedPane {
		case FocusDBTree:
			var treeModel tea.Model
			treeModel, cmd = m.dbtree.Update(msg)
			m.dbtree = treeModel.(dbtree.DBTreeModel)
			cmds = append(cmds, cmd)
		case FocusTableView:
			m.tableview, cmd = m.tableview.Update(msg)
			cmds = append(cmds, cmd)
		case FocusSQLCommandBar:
			var sqlModel tea.Model
			sqlModel, cmd = m.sqlCommandBar.Update(msg)
			m.sqlCommandBar = sqlModel.(sqlcommandbar.SQLCommandBarModel)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	} else {
		// Legacy mode: route to single model
		m.model, cmd = m.model.Update(msg)
		return m, cmd
	}
}

func (m rootScreenModel) View() string {
	var baseView string

	if m.useSplitLayout {
		baseView = m.renderSplitLayout()
	} else {
		baseView = m.model.View()
	}

	if m.notification == nil {
		return baseView
	}

	// Create the notification view
	style := m.notification.GetStyle()
	notifView := style.Render(m.notification.Message)

	if m.width > 0 && m.height > 0 {
		notifOverlay := lipgloss.Place(
			m.width,
			1,
			lipgloss.Right,
			lipgloss.Top,
			notifView,
		)

		return notifOverlay + baseView
	}

	return notifView + "\n" + baseView
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

	content := m.dbtree.View()
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

	content := m.tableview.View()
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

	content := m.sqlCommandBar.View()
	return borderStyle.Render(content)
}

// this is the switcher which will switch between screens
func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init() // must return .Init() to initialize the screen (and here the magic happens)
}
