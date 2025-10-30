package database

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5"
)

type Schema struct {
	Name string
}

func NewSchema(name string) *Schema {
	return &Schema{Name: name}
}

func (db *DatabaseConnection) ParseSchemas() ([]*Schema, error) {
	if !db.Connected {
		if err := db.Connect(); err != nil {
			return nil, err
		}
	}

	rows, err := db.connection.Query(context.Background(), "SELECT nspname from pg_namespace")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schemas, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Schema, error) {
		var name string
		err := row.Scan(&name)
		return NewSchema(name), err
	})
	// order "public" first
	sort.SliceStable(schemas, func(i, j int) bool {
		return schemas[i].Name == "public"
	})
	db.Schemas = schemas
	return schemas, err
}
