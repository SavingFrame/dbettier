package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type tableType int

const (
	baseTableType tableType = iota
	viewTableType
)

type Table struct {
	Name     string
	Type     tableType
	database *Database
}

func NewTable(name string, db *Database, tableType tableType) *Table {
	return &Table{Name: name, database: db, Type: tableType}
}

func (s *Schema) LoadTables() ([]*Table, error) {
	db := s.database
	if db.Connected {
		db.Connect()
	}
	q := "select table_name, table_type from information_schema.tables where table_schema = $1 ORDER BY table_type, table_name"

	rows, err := db.connection.Query(context.Background(), q, s.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Table, error) {
		var tableName, tableTypeRaw string
		var tableType tableType
		if err := row.Scan(&tableName, &tableTypeRaw); err != nil {
			return nil, err
		}

		switch tableTypeRaw {
		case "BASE TABLE":
			tableType = baseTableType
		case "VIEW":
			tableType = viewTableType
		}
		return NewTable(tableName, s.database, tableType), nil
	})
	return tables, err
}
