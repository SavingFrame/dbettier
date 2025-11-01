package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTables(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL
		)`,
		`CREATE TABLE test_schema.products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE test_schema.orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES test_schema.users(id)
		)`,
	)

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var testSchema *Schema
	for _, schema := range schemas {
		if schema.Name == "test_schema" {
			testSchema = schema
			break
		}
	}
	require.NotNil(t, testSchema)

	tables, err := testSchema.LoadTables()
	require.NoError(t, err, "Failed to load tables")
	assert.Len(t, tables, 3, "Expected 3 tables")

	tableNames := make([]string, len(tables))
	for i, table := range tables {
		tableNames[i] = table.Name

		assert.NotNil(t, table.Schema, "Table %s has nil Schema reference", table.Name)
		assert.Equal(t, testSchema, table.Schema, "Table %s references wrong schema", table.Name)
	}

	expectedTables := []string{"users", "products", "orders"}
	for _, expected := range expectedTables {
		assert.Contains(t, tableNames, expected, "Expected table %s not found", expected)
	}

	t.Logf("Successfully loaded %d tables", len(tables))
}

func TestLoadTablesEmpty(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "empty_schema")
	defer DropSchemas(t, db, "empty_schema")

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var emptySchema *Schema
	for _, schema := range schemas {
		if schema.Name == "empty_schema" {
			emptySchema = schema
			break
		}
	}
	require.NotNil(t, emptySchema)

	tables, err := emptySchema.LoadTables()
	require.NoError(t, err, "LoadTables should not error on empty schema")
	assert.Empty(t, tables, "Expected 0 tables in empty schema")
}

func TestLoadTablesTypes(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.base_table (
			id SERIAL PRIMARY KEY,
			value TEXT
		)`,
		`CREATE VIEW test_schema.test_view AS 
			SELECT id, value FROM test_schema.base_table`,
	)

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var testSchema *Schema
	for _, schema := range schemas {
		if schema.Name == "test_schema" {
			testSchema = schema
			break
		}
	}
	require.NotNil(t, testSchema)

	tables, err := testSchema.LoadTables()
	require.NoError(t, err)

	assert.Len(t, tables, 2, "Expected 2 items (table + view)")

	foundTable := false
	foundView := false
	for _, table := range tables {
		if table.Name == "base_table" {
			assert.Equal(t, "BASE TABLE", table.Type.String(), "base_table should be BASE TABLE type")
			foundTable = true
		}
		if table.Name == "test_view" {
			assert.Equal(t, "VIEW", table.Type.String(), "test_view should be VIEW type")
			foundView = true
		}
	}

	assert.True(t, foundTable, "Should have found base_table")
	assert.True(t, foundView, "Should have found test_view")
}

func TestLoadTablesSorting(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.zzz_table (id SERIAL PRIMARY KEY)`,
		`CREATE TABLE test_schema.aaa_table (id SERIAL PRIMARY KEY)`,
		`CREATE VIEW test_schema.mmm_view AS SELECT 1`,
	)

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var testSchema *Schema
	for _, schema := range schemas {
		if schema.Name == "test_schema" {
			testSchema = schema
			break
		}
	}
	require.NotNil(t, testSchema)

	tables, err := testSchema.LoadTables()
	require.NoError(t, err)
	require.Len(t, tables, 3)

	// Tables should be sorted by type, then name
	// So: all BASE TABLEs first (sorted by name), then VIEWs (sorted by name)

	assert.Equal(t, "aaa_table", tables[0].Name)
	assert.Equal(t, "BASE TABLE", tables[0].Type.String())

	assert.Equal(t, "zzz_table", tables[1].Name)
	assert.Equal(t, "BASE TABLE", tables[1].Type.String())

	assert.Equal(t, "mmm_view", tables[2].Name)
	assert.Equal(t, "VIEW", tables[2].Type.String())
}

func TestFindTable(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.users (id SERIAL PRIMARY KEY)`,
		`CREATE TABLE test_schema.products (id SERIAL PRIMARY KEY)`,
	)

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var testSchema *Schema
	for _, schema := range schemas {
		if schema.Name == "test_schema" {
			testSchema = schema
			break
		}
	}
	require.NotNil(t, testSchema)

	_, err = testSchema.LoadTables()
	require.NoError(t, err)

	usersTable := testSchema.FindTable("users")
	require.NotNil(t, usersTable, "FindTable should find existing table")
	assert.Equal(t, "users", usersTable.Name)
	assert.Equal(t, testSchema, usersTable.Schema)

	nonExistent := testSchema.FindTable("non_existent_table")
	assert.Nil(t, nonExistent, "FindTable should return nil for non-existent table")

	upperCase := testSchema.FindTable("USERS")
	assert.Nil(t, upperCase, "FindTable should be case-sensitive")
}
