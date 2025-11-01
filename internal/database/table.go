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

func (t tableType) String() string {
	switch t {
	case baseTableType:
		return "BASE TABLE"
	case viewTableType:
		return "VIEW"
	default:
		return "UNKNOWN"
	}
}

type Table struct {
	Name    string
	Type    tableType
	Schema  *Schema
	Columns []*Column
}

func NewTable(name string, schema *Schema, tableType tableType) *Table {
	return &Table{Name: name, Schema: schema, Type: tableType}
}

func (s *Schema) LoadTables() ([]*Table, error) {
	db := s.Database
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
		return NewTable(tableName, s, tableType), nil
	})
	s.Tables = tables
	return tables, err
}

func (s *Schema) FindTable(name string) *Table {
	for _, table := range s.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}
