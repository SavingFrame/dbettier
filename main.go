package main

import (
	"fmt"
	"os"

	"github.com/SavingFrame/dbettier/internal/components"
	"github.com/SavingFrame/dbettier/internal/database"
	tea "github.com/charmbracelet/bubbletea"
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

	database.LoadConnections()
	if _, err := tea.NewProgram(components.RootScreen(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
