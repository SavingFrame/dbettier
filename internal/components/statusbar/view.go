package statusbar

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/SavingFrame/dbettier/internal/theme"
)

func statusNuggetStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Base).
		Padding(0, 1)
}

func statusBarStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Foreground(colors.Subtle).
		Background(colors.Surface)
}

func statusStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return lipgloss.NewStyle().
		Inherit(statusBarStyle()).
		Foreground(colors.Base).
		Background(colors.Primary).
		Padding(0, 1).
		MarginRight(1)
}

func encodingStyle() lipgloss.Style {
	colors := theme.Current().Colors
	return statusNuggetStyle().
		Background(colors.Secondary).
		Align(lipgloss.Right)
}

func statusTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Inherit(statusBarStyle())
}

func (s StatusBarModel) RenderContent() string {
	encoding := encodingStyle().Render("UTF-8")
	status := statusStyle().Render(s.editorStatus)
	mode := statusTextStyle().Render("NORMAL")
	bar := lipgloss.JoinHorizontal(lipgloss.Top, encoding, status, mode)

	return statusBarStyle().Width(s.width).Height(s.height).Render(bar)
}

func (s StatusBarModel) View() tea.View {
	var v tea.View
	v.SetContent(s.RenderContent())
	return v
}

func (s StatusBarModel) Init() tea.Cmd {
	return nil
}
