package database

type Table struct {
	Name     string
	database *Database
}

func NewTable(name string, db *Database) *Table {
	return &Table{Name: name, database: db}
}
