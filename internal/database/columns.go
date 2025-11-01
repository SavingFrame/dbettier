package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5"
)

type Column struct {
	table           *Table
	Name            string
	ColumnDefault   sql.NullString
	Nullable        bool
	DataType        string
	UserDefinedType string
	MaxLength       sql.NullInt32
	IsPrimaryKey    bool
}

func NewColumn(
	name string,
	table *Table,
	columnDefault sql.NullString,
	nullable bool,
	dataType string,
	userDefinedType string,
	maxLength sql.NullInt32,
	isPrimaryKey bool,
) *Column {
	return &Column{
		Name:            name,
		ColumnDefault:   columnDefault,
		Nullable:        nullable,
		DataType:        dataType,
		UserDefinedType: userDefinedType,
		MaxLength:       maxLength,
		IsPrimaryKey:    isPrimaryKey,
		table:           table,
	}
}

func (t *Table) LoadColumnsForTable() ([]*Column, error) {
	db := t.Schema.Database
	if !db.Connected {
		if err := db.Connect(); err != nil {
			return nil, err
		}
	}
	q := `SELECT column_name, is_nullable, data_type, character_maximum_length, udt_name, is_identity, column_default
FROM information_schema.columns
WHERE table_schema = $1
 AND table_name = $2
	ORDER BY ordinal_position`
	rows, err := db.connection.Query(context.Background(), q, t.Schema.Name, t.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Column, error) {
		var col Column
		var isNullable string
		var isPrimaryKey string
		var maxLength *int32
		var columnDefault *string
		if err := row.Scan(
			&col.Name,
			&isNullable,
			&col.DataType,
			&maxLength,
			&col.UserDefinedType,
			&isPrimaryKey,
			&columnDefault,
		); err != nil {
			return nil, err
		}

		// Convert string values to proper types
		col.Nullable = isNullable == "YES"
		col.IsPrimaryKey = isPrimaryKey == "YES"

		if maxLength != nil {
			col.MaxLength = sql.NullInt32{Int32: *maxLength, Valid: true}
		} else {
			col.MaxLength = sql.NullInt32{Valid: false}
		}

		if columnDefault != nil {
			col.ColumnDefault = sql.NullString{String: *columnDefault, Valid: true}
		} else {
			col.ColumnDefault = sql.NullString{Valid: false}
		}

		col.table = t

		return &col, nil
	})
	return columns, err
}

func (s *Schema) LoadColumns() (map[*Table][]*Column, error) {
	db := s.Database
	if !db.Connected {
		if err := db.Connect(); err != nil {
			return nil, err
		}
	}

	columnsByTable := make(map[*Table][]*Column)

	q := `SELECT column_name, is_nullable, data_type, character_maximum_length, udt_name, is_identity, column_default, table_name
  FROM information_schema.columns
 WHERE table_schema = $1
	ORDER BY table_name, ordinal_position`
	rows, err := db.connection.Query(context.Background(), q, s.Name)
	if err != nil {
		log.Printf("Error querying columns for schema %s: %v", s.Name, err)
		return nil, err
	}
	defer rows.Close()

	var table *Table
	_, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Column, error) {
		var col Column
		var tableName string
		var isNullable string
		var isPrimaryKey string
		var maxLength *int32
		var columnDefault *string
		if err := row.Scan(
			&col.Name,
			&isNullable,
			&col.DataType,
			&maxLength,
			&col.UserDefinedType,
			&isPrimaryKey,
			&columnDefault,
			&tableName,
		); err != nil {
			log.Printf("Error scanning row for column %s in table %s: %v", col.Name, tableName, err)
			return nil, err
		}
		if table != nil && table.Schema == s && table.Name == tableName {
			col.table = table
		} else {
			table = s.FindTable(tableName)
			col.table = table
		}
		if maxLength != nil {
			col.MaxLength = sql.NullInt32{Int32: *maxLength, Valid: true}
		} else {
			col.MaxLength = sql.NullInt32{Valid: false}
		}
		if columnDefault != nil {
			col.ColumnDefault = sql.NullString{String: *columnDefault, Valid: true}
		} else {
			col.ColumnDefault = sql.NullString{Valid: false}
		}
		col.Nullable = isNullable == "YES"
		col.IsPrimaryKey = isPrimaryKey == "YES"
		columnsByTable[table] = append(columnsByTable[table], &col)

		return &col, nil
	})
	if err != nil {
		log.Printf("Error collecting rows for columns in schema %s: %v", s.Name, err)
	}
	return columnsByTable, err
}
