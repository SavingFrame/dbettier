package components

import (
	"time"

	"github.com/SavingFrame/dbettier/internal/components/dbtree"
	"github.com/SavingFrame/dbettier/internal/components/notifications"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type rootScreenModel struct {
	model        tea.Model // this will hold the current screen model
	notification *notifications.Notification
	width        int
	height       int
	registry     *database.DBRegistry
}

func RootScreen(registry *database.DBRegistry) rootScreenModel {
	var rootModel tea.Model

	// sample conditional logic to start with a specific screen
	// notice that the screen methods Update and View have been modified
	// to accept a pointer *screenXModel instead of screenXModel
	// this will allow us to modify the model's state in the View method
	// if needed

	if registry.Count() > 0 {
		dbtreeScreen := dbtree.DBTreeScreen(registry)
		rootModel = &dbtreeScreen
	} else {
		screenOne := DBCreatorScreen(registry)
		rootModel = &screenOne
	}
	// } else {
	//     screen_two := screenTwo()
	//     rootModel = &screen_two
	// }

	return rootScreenModel{
		model:    rootModel,
		registry: registry,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init() // rest methods are just wrappers for the model's methods
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m rootScreenModel) View() string {
	baseView := m.model.View()

	if m.notification == nil {
		return baseView
	}

	// Create the notification view
	style := m.notification.GetStyle()
	notifView := style.Render(m.notification.Message)

	// If we have terminal dimensions, use lipgloss.Place for proper overlay
	if m.width > 0 && m.height > 0 {
		// Place notification in top-right corner
		notifOverlay := lipgloss.Place(
			m.width,
			1, // Only take 1 line at the top
			lipgloss.Right,
			lipgloss.Top,
			notifView,
		)

		// return baseView + "\n\n" + notifOverlay
		return notifOverlay + baseView
		// Split base view into lines
		// lines := strings.Split(baseView, "\n")
		// if len(lines) > 0 {
		// 	// Replace or overlay the first line with notification
		// 	lines[0] = notifOverlay
		// 	return strings.Join(lines, "\n")
		// }
	}

	// Fallback: just put notification at the top
	return notifView + "\n" + baseView
}

// this is the switcher which will switch between screens
func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init() // must return .Init() to initialize the screen (and here the magic happens)
}
