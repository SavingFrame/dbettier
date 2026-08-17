package sqlcommandbar

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/SavingFrame/dbettier/internal/database"
	"github.com/SavingFrame/dbettier/pkgs/editor"
)

func createTempDir() (string, error) {
	dir, err := os.MkdirTemp("", "dbettier-sql-")
	if err != nil {
		return "", err
	}
	return dir, err
}

func createTempContentFile(dir string, editor *editor.SQLEditor) (string, error) {
	name := filepath.Join(dir, "query.sql")
	content := editor.GetContent()
	if err := os.WriteFile(name, []byte(content), 0o700); err != nil {
		return "", nil
	}
	return name, nil
}

func createTempCredentialsFile(dir string, db database.Database) (string, error) {
	name := filepath.Join(dir, "postgres-language-server.jsonc")
	settings := map[string]any{
		"db": map[string]any{
			"host":            db.Host,
			"port":            db.Port,
			"username":        db.Username,
			"password":        db.Password,
			"database":        db.Database,
			"connTimeoutSecs": 10,
		},
	}

	content, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", err
	}
	content = append(content, '\n')

	if err := os.WriteFile(name, content, 0o600); err != nil {
		return "", err
	}
	return name, nil
}

func updateEditorFromFile(name string, editor *editor.SQLEditor) error {
	bytes, err := os.ReadFile(name)
	if err != nil {
		return nil
	}
	editor.SetContent(string(bytes[:]))
	return nil
}

func deleteTempBufferFile(name string) {
	os.Remove(name)
}
