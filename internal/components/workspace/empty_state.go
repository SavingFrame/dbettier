package workspace

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/SavingFrame/dbettier/internal/theme"
)

func (w *Workspace) renderNoTabTableState() string {
	title := emptyWorkspaceTitleStyle().Render("󰆍 Welcome to dbettier")
	subtitle := emptyWorkspaceSubtitleStyle().Render("No tabs open yet. Pick where you want to start.")

	sep := emptyWorkspaceHintStyle().Render(" ")
	hint1 := emptyWorkspaceBulletStyle().Render("•") + sep + emptyWorkspaceHintStyle().Render("Press Enter on a table in the left tree")
	hint2 := emptyWorkspaceBulletStyle().Render("•") + sep + emptyWorkspaceHintStyle().Render("Press c on a selected table to open a query tab")
	hint3 := emptyWorkspaceBulletStyle().Render("•") + sep + emptyWorkspaceHintStyle().Render("Run SQL from the editor with Alt+Enter")

	content := strings.Join([]string{
		title,
		subtitle,
		"",
		hint1,
		hint2,
		hint3,
	}, "\n")

	contentWidth := lipgloss.Width(content)
	surfaceFill := lipgloss.NewStyle().Background(theme.Current().Colors.Surface)
	content = lipgloss.PlaceHorizontal(
		contentWidth,
		lipgloss.Left,
		content,
		lipgloss.WithWhitespaceStyle(surfaceFill),
	)

	card := emptyWorkspaceCardStyle().Render(content)
	return centerWorkspaceBlock(w.TableViewSize.width, w.TableViewSize.height, card)
}

func (w *Workspace) renderNoTabSQLState() string {
	title := emptyWorkspaceTitleStyle().Render("󰆍 SQL editor is waiting")
	subtitle := emptyWorkspaceSubtitleStyle().Render("Open a tab from the tree to start writing queries.")

	content := strings.Join([]string{title, subtitle}, "\n")
	contentWidth := lipgloss.Width(content)
	surfaceFill := lipgloss.NewStyle().Background(theme.Current().Colors.Surface)
	content = lipgloss.PlaceHorizontal(
		contentWidth,
		lipgloss.Left,
		content,
		lipgloss.WithWhitespaceStyle(surfaceFill),
	)

	card := emptyWorkspaceCardStyle().Render(content)
	return centerWorkspaceBlock(w.SQLCommandBarSize.width, w.SQLCommandBarSize.height, card)
}

func centerWorkspaceBlock(width, height int, content string) string {
	if width <= 0 || height <= 0 {
		return content
	}

	fill := lipgloss.NewStyle().Background(theme.Current().Colors.Base)
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceStyle(fill),
	)
}
