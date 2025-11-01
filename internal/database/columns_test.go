package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadColumns(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT true
		)`,
		`CREATE TABLE test_schema.products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			price NUMERIC(10, 2),
			description TEXT,
			stock_quantity INTEGER DEFAULT 0
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

	// Load tables (required before LoadColumns)
	_, err = testSchema.LoadTables()
	require.NoError(t, err)

	// Test LoadColumns
	columns, err := testSchema.LoadColumns()
	require.NoError(t, err, "Failed to load columns")
	assert.NotEmpty(t, columns, "No columns loaded")

	// Verify columns for each table
	expectedColumns := map[string][]string{
		"users":    {"id", "username", "email", "created_at", "is_active"},
		"products": {"id", "name", "price", "description", "stock_quantity"},
	}

	for table, cols := range columns {
		expected, ok := expectedColumns[table.Name]
		if !ok {
			continue // Skip non-test tables
		}

		columnNames := make([]string, len(cols))
		for i, col := range cols {
			columnNames[i] = col.Name
		}

		for _, expectedCol := range expected {
			assert.Contains(t, columnNames, expectedCol,
				"Expected column %s.%s not found", table.Name, expectedCol)
		}

		assert.Len(t, cols, len(expected),
			"Table %s: expected %d columns, got %d", table.Name, len(expected), len(cols))
	}
}

func TestLoadColumnsEmpty(t *testing.T) {
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

	// Load tables (should be empty)
	_, err = emptySchema.LoadTables()
	require.NoError(t, err)

	// Load columns (should return empty map)
	columns, err := emptySchema.LoadColumns()
	require.NoError(t, err, "LoadColumns should not error on empty schema")
	assert.Empty(t, columns, "Expected 0 columns in empty schema")
}

func TestLoadColumnsTableReference(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100)
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

	// Load tables
	_, err = testSchema.LoadTables()
	require.NoError(t, err)

	// Load columns
	columns, err := testSchema.LoadColumns()
	require.NoError(t, err)

	// Verify that each column has a reference to its table
	for table, cols := range columns {
		assert.NotNil(t, table, "Found nil table in columns map")

		for _, col := range cols {
			assert.NotNil(t, col.table,
				"Column %s.%s has nil table reference", table.Name, col.Name)
			assert.Equal(t, table, col.table,
				"Column %s.%s references wrong table", table.Name, col.Name)
		}
	}
}

func TestLoadColumnsForTable(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100),
			age INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	require.NoError(t, err)
	require.Len(t, tables, 1)

	usersTable := tables[0]

	// Test LoadColumnsForTable
	columns, err := usersTable.LoadColumnsForTable()
	require.NoError(t, err, "LoadColumnsForTable should not error")
	assert.Len(t, columns, 5, "Expected 5 columns")

	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = col.Name

		// Verify each column has reference to the table
		assert.NotNil(t, col.table, "Column %s should have table reference", col.Name)
		assert.Equal(t, usersTable, col.table, "Column %s has wrong table reference", col.Name)
	}

	expectedColumns := []string{"id", "username", "email", "age", "created_at"}
	for _, expected := range expectedColumns {
		assert.Contains(t, columnNames, expected, "Expected column %s not found", expected)
	}
}

func TestColumnDataTypes(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.data_types (
			int_col INTEGER,
			bigint_col BIGINT,
			varchar_col VARCHAR(100),
			text_col TEXT,
			bool_col BOOLEAN,
			numeric_col NUMERIC(10, 2),
			date_col DATE,
			timestamp_col TIMESTAMP,
			json_col JSONB
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
	require.NoError(t, err)
	require.Len(t, tables, 1)

	columns, err := tables[0].LoadColumnsForTable()
	require.NoError(t, err)

	// Verify data types
	expectedTypes := map[string]string{
		"int_col":       "integer",
		"bigint_col":    "bigint",
		"varchar_col":   "character varying",
		"text_col":      "text",
		"bool_col":      "boolean",
		"numeric_col":   "numeric",
		"date_col":      "date",
		"timestamp_col": "timestamp without time zone",
		"json_col":      "jsonb",
	}

	for _, col := range columns {
		expectedType, ok := expectedTypes[col.Name]
		if ok {
			assert.Equal(t, expectedType, col.DataType,
				"Column %s has wrong data type", col.Name)
		}
	}
}

func TestColumnNullability(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.nullability (
			required_col VARCHAR(50) NOT NULL,
			optional_col VARCHAR(50)
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
	require.NoError(t, err)

	columns, err := tables[0].LoadColumnsForTable()
	require.NoError(t, err)

	for _, col := range columns {
		if col.Name == "required_col" {
			assert.False(t, col.Nullable, "required_col should not be nullable")
		}
		if col.Name == "optional_col" {
			assert.True(t, col.Nullable, "optional_col should be nullable")
		}
	}
}

func TestColumnDefaults(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	ExecQueries(t, db,
		`CREATE TABLE test_schema.defaults (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50),
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	require.NoError(t, err)

	columns, err := tables[0].LoadColumnsForTable()
	require.NoError(t, err)

	for _, col := range columns {
		if col.Name == "id" {
			// SERIAL creates a default nextval()
			assert.True(t, col.ColumnDefault.Valid,
				"id column should have a default value")
		}
		if col.Name == "name" {
			assert.False(t, col.ColumnDefault.Valid,
				"name column should not have a default value")
		}
		if col.Name == "active" {
			assert.True(t, col.ColumnDefault.Valid,
				"active column should have a default value")
		}
		if col.Name == "created_at" {
			assert.True(t, col.ColumnDefault.Valid,
				"created_at column should have a default value")
		}
	}
}
