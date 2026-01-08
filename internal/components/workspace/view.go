package workspace

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// RenderTabBar renders the tab bar
func (w *Workspace) RenderTabBar() string {
	if len(w.tabs) == 0 {
		return ""
	}

	activeStyle := activeTabStyle()
	inactiveStyle := inactiveTabStyle()
	gapStyle := tabGapStyle()

	var renderedTabs []string
	totalWidth := 0

	// Calculate which tabs to show based on scroll offset and available width
	visibleTabs := w.calculateVisibleTabs()

	for _, idx := range visibleTabs {
		tab := w.tabs[idx]
		isActive := idx == w.activeIndex

		// Build tab content: icon + name + close button
		icon := iconStyle(tab.Type).Render(tab.Icon())
		name := tab.Name
		closeBtn := closeButtonStyle(isActive).Render("×")

		content := fmt.Sprintf("%s%s %s", icon, name, closeBtn)

		var tabView string
		if isActive {
			tabView = activeStyle.Render(content)
		} else {
			tabView = inactiveStyle.Render(content)
		}

		// Wrap with zone for mouse interaction
		tabZone := zone.Mark(fmt.Sprintf("tab-%d", idx), tabView)
		closeZone := zone.Mark(fmt.Sprintf("tab-close-%d", idx), "")

		// We need to embed the close zone within the tab
		// For simplicity, we'll handle click detection in update based on position
		_ = closeZone

		renderedTabs = append(renderedTabs, tabZone)
		totalWidth += lipgloss.Width(tabView)
	}

	// Join all tabs horizontally
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Add gap filler to extend the bottom border to the full width
	rowWidth := lipgloss.Width(row)
	if rowWidth < w.width && w.width > 0 {
		gapWidth := w.width - rowWidth - 2 // -2 for some margin
		if gapWidth > 0 {
			gap := gapStyle.Render(strings.Repeat(" ", gapWidth))
			row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		}
	}

	return row
}

// calculateVisibleTabs returns the indices of tabs that should be visible
func (w *Workspace) calculateVisibleTabs() []int {
	if w.width <= 0 {
		// Return all tabs if width not set yet
		indices := make([]int, len(w.tabs))
		for i := range w.tabs {
			indices[i] = i
		}
		return indices
	}

	// Estimate tab width (icon + name + close button + padding + borders)
	// This is approximate; actual width depends on tab name length
	estimateTabWidth := func(tab Tab) int {
		// icon (2) + space (1) + name + space (1) + close (1) + padding (2) + borders (2)
		return len(tab.Name) + 9
	}

	var visibleIndices []int
	currentWidth := 0
	availableWidth := w.width - 4 // Leave some margin

	// First, ensure active tab is included by adjusting scroll offset
	// Calculate total width of tabs before active tab
	widthBeforeActive := 0
	for i := 0; i < w.activeIndex; i++ {
		widthBeforeActive += estimateTabWidth(w.tabs[i])
	}

	// If active tab would be scrolled out of view, adjust scroll offset
	if w.scrollOffset > w.activeIndex {
		w.scrollOffset = w.activeIndex
	}

	// Calculate width from scroll offset to active tab
	widthToActive := 0
	for i := w.scrollOffset; i <= w.activeIndex && i < len(w.tabs); i++ {
		widthToActive += estimateTabWidth(w.tabs[i])
	}

	// If active tab is beyond visible area, increase scroll offset
	for widthToActive > availableWidth && w.scrollOffset < w.activeIndex {
		widthToActive -= estimateTabWidth(w.tabs[w.scrollOffset])
		w.scrollOffset++
	}

	// Now collect visible tabs starting from scroll offset
	for i := w.scrollOffset; i < len(w.tabs); i++ {
		tabWidth := estimateTabWidth(w.tabs[i])
		if currentWidth+tabWidth > availableWidth && len(visibleIndices) > 0 {
			// Don't break if this is the active tab - always show it
			if i != w.activeIndex {
				break
			}
		}
		visibleIndices = append(visibleIndices, i)
		currentWidth += tabWidth
	}

	return visibleIndices
}

// GetTabIndexAtPosition returns the tab index at a given x position, or -1 if none
func (w *Workspace) GetTabIndexAtPosition(x int) int {
	if len(w.tabs) == 0 {
		return -1
	}

	activeStyle := activeTabStyle()
	inactiveStyle := inactiveTabStyle()

	currentX := 0
	visibleTabs := w.calculateVisibleTabs()

	for _, idx := range visibleTabs {
		tab := w.tabs[idx]
		isActive := idx == w.activeIndex

		// Calculate tab width
		icon := iconStyle(tab.Type).Render(tab.Icon())
		name := tab.Name
		closeBtn := closeButtonStyle(isActive).Render("×")
		content := fmt.Sprintf("%s%s %s", icon, name, closeBtn)

		var tabWidth int
		if isActive {
			tabWidth = lipgloss.Width(activeStyle.Render(content))
		} else {
			tabWidth = lipgloss.Width(inactiveStyle.Render(content))
		}

		if x >= currentX && x < currentX+tabWidth {
			return idx
		}
		currentX += tabWidth
	}

	return -1
}

// IsCloseButtonClick checks if a click at position x within a tab is on the close button
func (w *Workspace) IsCloseButtonClick(tabIndex, relativeX int) bool {
	if tabIndex < 0 || tabIndex >= len(w.tabs) {
		return false
	}

	tab := w.tabs[tabIndex]
	isActive := tabIndex == w.activeIndex

	// Calculate tab content width
	icon := iconStyle(tab.Type).Render(tab.Icon())
	name := tab.Name
	closeBtn := closeButtonStyle(isActive).Render("×")
	content := fmt.Sprintf("%s%s %s", icon, name, closeBtn)

	var style lipgloss.Style
	if isActive {
		style = activeTabStyle()
	} else {
		style = inactiveTabStyle()
	}

	tabWidth := lipgloss.Width(style.Render(content))

	// Close button is at the right side of the tab (last ~3 characters including padding)
	closeButtonStart := tabWidth - 4
	return relativeX >= closeButtonStart
}

func (w Workspace) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(w.RenderTabBar())
	return v
}
