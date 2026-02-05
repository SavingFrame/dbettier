package editor

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/database"
)

type EditorMode int

const (
	EditorModeNormal EditorMode = iota
	EditorModeInsert
)

type SQLEditor struct {
	viewport viewport.Model
	mode     EditorMode
	buffer   *buffer
	cursor   *editorCursor
	readonly bool

	registry *database.DBRegistry
	ready    bool
}

func NewEditorModel(lines []string, readonly bool) SQLEditor {
	if lines == nil {
		lines = []string{"This is a simple text editor.", "", "Start typing..."}
	}

	m := SQLEditor{
		viewport: viewport.New(),
		mode:     EditorModeNormal,
		buffer:   &buffer{lines: lines},
		cursor:   newEditorCursor(0, 0),
		readonly: readonly,
	}

	return m
}

func (m SQLEditor) Init() tea.Cmd {
	return nil
}

func (m *SQLEditor) SetSize(width, height int) {
	if !m.ready {
		m.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
		m.viewport.SetContent(strings.Join(m.buffer.lines, "\n"))
		m.ready = true
	}
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
}

func (m *SQLEditor) Focus() tea.Cmd {
	return m.cursor.virtualCursor.Focus()
}

func (m *SQLEditor) Blur() tea.Cmd {
	m.cursor.virtualCursor.Blur()
	return nil
}

func (m *SQLEditor) SetContent(c string) {
	c = strings.Trim(c, " ")
	contentLines := strings.Split(c, "\n")
	m.buffer.lines = contentLines
	m.cursor.moveLastSymbol(m.buffer.lines)
}

func (m *SQLEditor) GetContent() string {
	return strings.Join(m.buffer.lines, "\n")
}

func (m *SQLEditor) Mode() EditorMode {
	return m.mode
}

func (m *SQLEditor) IsNormalMode() bool {
	return m.mode == EditorModeNormal
}

func (m *SQLEditor) SetReadonly(v bool) {
	m.readonly = v
}

func (m *SQLEditor) IsReadonly() bool {
	return m.readonly
}
