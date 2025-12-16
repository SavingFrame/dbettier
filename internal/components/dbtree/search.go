package dbtree

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// TreeSearchMatch represents a match in the tree.
// Path follows same format as treeCursor.path
type TreeSearchMatch struct {
	Path []int
	Name string // The matched name (for display)
}

type TreeSearch struct {
	mode       bool              // Whether search input is active
	query      string            // Current search query
	matches    []TreeSearchMatch // All matching nodes
	matchIndex int               // Current match index (-1 if no matches)
}

// IsMatch checks if a path matches any search result.
// Returns (isMatch, isActive)
func (s *TreeSearch) IsMatch(path []int) (bool, bool) {
	if len(s.matches) == 0 {
		return false, false
	}

	for i, match := range s.matches {
		if pathsEqual(path, match.Path) {
			return true, i == s.matchIndex
		}
	}
	return false, false
}

func (s *TreeSearch) Enable() {
	s.mode = true
	s.query = ""
	s.matches = nil
	s.matchIndex = -1
}

// Clear clears the search state.
func (s *TreeSearch) Clear() {
	s.mode = false
	s.query = ""
	s.matches = nil
	s.matchIndex = -1
}

// HandleInput handles key input during search mode.
func (s *TreeSearch) HandleInput(msg tea.KeyMsg, tree *TreeState, cursor *TreeCursor) {
	switch msg.String() {
	case "esc":
		// Exit search mode and clear highlights
		s.Clear()

	case "enter":
		// Confirm search and exit search mode
		s.mode = false
		if len(s.matches) > 0 && s.matchIndex >= 0 {
			// Jump to current match
			match := s.matches[s.matchIndex]
			tree.cursor.SetPath(match.Path)
			// Should we move adjusting
		}

	case "backspace":
		if len(s.query) > 0 {
			s.query = s.query[:len(s.query)-1]
			s.UpdateMatches(tree, cursor)
		}

	default:
		// Add character to search query (only printable characters)
		if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] < 127 {
			s.query += msg.String()
			s.UpdateMatches(tree, cursor)
		}
	}
}

// UpdateMatches finds all tree nodes matching the search query.
func (s *TreeSearch) UpdateMatches(tree *TreeState, cursor *TreeCursor) {
	s.matches = nil
	s.matchIndex = -1

	if s.query == "" {
		return
	}

	query := strings.ToLower(s.query)

	// Search through all visible nodes in the tree
	for dbIdx, db := range tree.databases {
		// Check database name
		if strings.Contains(strings.ToLower(db.name), query) {
			s.matches = append(s.matches, TreeSearchMatch{
				Path: []int{dbIdx},
				Name: db.name,
			})
		}

		// Search in schemas (only if database is expanded)
		if db.expanded {
			for schemaIdx, schema := range db.schemas {
				// Check schema name
				if strings.Contains(strings.ToLower(schema.name), query) {
					s.matches = append(s.matches, TreeSearchMatch{
						Path: []int{dbIdx, schemaIdx},
						Name: schema.name,
					})
				}

				// Search in tables (only if schema is expanded)
				if schema.expanded {
					for tableIdx, table := range schema.tables {
						// Check table name
						if strings.Contains(strings.ToLower(table.name), query) {
							s.matches = append(s.matches, TreeSearchMatch{
								Path: []int{dbIdx, schemaIdx, tableIdx},
								Name: table.name,
							})
						}

						// Search in columns (only if table is expanded)
						if table.expanded {
							for colIdx, col := range table.columns {
								// Check column name
								if strings.Contains(strings.ToLower(col.name), query) {
									s.matches = append(s.matches, TreeSearchMatch{
										Path: []int{dbIdx, schemaIdx, tableIdx, colIdx},
										Name: col.name,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	// If we have matches, set index to first match and jump to it
	if len(s.matches) > 0 {
		s.matchIndex = 0
		match := s.matches[0]
		cursor.SetPath(match.Path)
	}
}

// NextMatch moves to the next search match.
func (s *TreeSearch) NextMatch(cursor *TreeCursor) {
	if len(s.matches) == 0 {
		return
	}

	s.matchIndex = (s.matchIndex + 1) % len(s.matches)
	cursor.SetPath(s.matches[s.matchIndex].Path)
}

// PrevMatch moves to the previous search match.
func (s *TreeSearch) PrevMatch(cursor *TreeCursor) {
	if len(s.matches) == 0 {
		return
	}

	s.matchIndex--
	if s.matchIndex < 0 {
		s.matchIndex = len(s.matches) - 1
	}
	cursor.SetPath(s.matches[s.matchIndex].Path)
}

func (s *TreeSearch) IsActive() bool {
	return s.mode
}

func (s *TreeSearch) Matches() []TreeSearchMatch {
	return s.matches
}

func (s *TreeSearch) Query() string {
	return s.query
}

func (s *TreeSearch) CurrentMatchIndex() int {
	return s.matchIndex
}

// pathsEqual compares two paths for equality.
func pathsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
