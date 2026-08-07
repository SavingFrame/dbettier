package sqlcommandbar

import (
	"os"

	"github.com/SavingFrame/dbettier/pkgs/editor"
)

func createTempBufferFile(editor *editor.SQLEditor) (string, error) {
	content := editor.GetContent()
	f, err := os.CreateTemp("", "dbettier-sql-*.sql")
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
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
