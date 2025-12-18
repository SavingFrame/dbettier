// Package main is the entry point for dbettier, a terminal-based database management tool.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/SavingFrame/dbettier/internal/components"
	"github.com/SavingFrame/dbettier/internal/database"
	zone "github.com/lrstanley/bubblezone/v2"
)

// setupDebugLog initializes debug logging if DEBUG env var is set.
// Usage: Just set DEBUG=1 when running the app: DEBUG=1 go run main.go
func setupDebugLog() func() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		return func() { f.Close() }
	}
	return func() {}
}

func main() {
	cleanup := setupDebugLog()
	defer cleanup()
	zone.NewGlobal()

	// Create database registry and load connections
	registry := database.NewDBRegistry()
	if err := registry.LoadFromFile(".connections.json"); err != nil {
		fmt.Println("Warning: could not load connections:", err)
	}

	v := components.RootScreen(registry)

	if _, err := tea.NewProgram(v).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
