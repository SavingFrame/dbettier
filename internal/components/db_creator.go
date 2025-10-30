package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton     = focusedStyle.Render("[ Submit ]")
	testButton        = fmt.Sprintf("[ %s ]", blurredStyle.Render("Test Connection"))
	focusedTestButton = focusedStyle.Render("[ Test Connection ]")
	blurredButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	countButtons      = 2
)

type DBCreatorModel struct {
	focusIndex   int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	dbTestStatus string
	err          string
	registry     *database.DBRegistry
}

func DBCreatorScreen(registry *database.DBRegistry) DBCreatorModel {
	m := DBCreatorModel{
		inputs:   make([]textinput.Model, 5),
		registry: registry,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Width = 20

		switch i {
		case 0:
			t.Placeholder = "Host"
			t.SetValue("localhost")
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.SetValue("5432")
			t.Placeholder = "Port"
			t.CharLimit = 5
		case 2:
			t.SetValue("postgres")
			t.Placeholder = "Username"
			t.CharLimit = 64
		case 3:
			t.SetValue("password")
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 4:
			t.Placeholder = "Database"
			t.Focus()
			m.focusIndex = 4
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func (m DBCreatorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DBCreatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case testDatabaseResult:
		m.dbTestStatus = string(msg)
		return m, nil
	case createDatabaseResult:
		if bool(msg) {
			m.dbTestStatus = "Database connection created successfully!"
		} else {
			m.dbTestStatus = "Failed to create database connection."
		}
		return m, nil

	case errMsg:
		m.err = msg.Error()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			totalFocusable := len(m.inputs) + countButtons - 1
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == 5 {
				port := m.inputs[1].Value()
				portInt, err := strconv.Atoi(port)
				if err != nil {
					fmt.Println("Invalid port number")
				}
				return m, createDatabase(m.inputs[0].Value(), m.inputs[2].Value(), m.inputs[3].Value(), portInt, m.inputs[4].Value(), m.registry)
			}
			if s == "enter" && m.focusIndex == 6 { // Test Connection button
				port := m.inputs[1].Value()
				portInt, err := strconv.Atoi(port)
				if err != nil {
					fmt.Println("Invalid port number")
				}
				m.dbTestStatus = "Testing connection..."
				return m, testDatabase(m.inputs[0].Value(), m.inputs[2].Value(), m.inputs[3].Value(), portInt, m.inputs[4].Value())
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > totalFocusable {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = totalFocusable
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *DBCreatorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type (
	errMsg             struct{ err error }
	testDatabaseResult string
)

func (e errMsg) Error() string { return e.err.Error() }
func testDatabase(host, username, password string, port int, db string) tea.Cmd {
	return func() tea.Msg {
		conn := database.NewDatabase(host, username, password, port, db)
		_, result := conn.Test()
		return testDatabaseResult(result)
	}
}

type createDatabaseResult bool

func createDatabase(host, username, password string, port int, db string, registry *database.DBRegistry) tea.Cmd {
	return func() tea.Msg {
		conn := database.NewDatabase(host, username, password, port, db)
		err := conn.SaveAndConnect(registry, ".connections.json")
		if err != nil {
			return errMsg{err}
		}
		return createDatabaseResult(err == nil)
	}
}

func (m DBCreatorModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	tButton := &testButton
	switch m.focusIndex {
	case 5:
		button = &focusedButton
	case 6:
		tButton = &focusedTestButton
	}
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s%s\n\n", *button, *tButton)

	if m.dbTestStatus != "" {
		fmt.Fprintf(&b, "\n%s\n\n", successStyle.Render(m.dbTestStatus))
	}
	if m.err != "" {
		fmt.Fprintf(&b, "\n%s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.err))
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
