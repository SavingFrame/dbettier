package workspace

import sharedcomponents "github.com/SavingFrame/dbettier/internal/components/shared_components"

type OpenQueryTabMsg struct {
	Query      sharedcomponents.QueryCompiler
	DatabaseID string
}
