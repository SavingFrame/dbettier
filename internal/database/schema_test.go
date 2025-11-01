package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain manages the lifecycle of the shared test container.
// It runs once before all tests and cleans up after all tests complete.
func TestMain(m *testing.M) {
	code := m.Run()

	TerminateSharedContainer()

	os.Exit(code)
}

func TestParseSchemas(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "schema_a")
	CreateSchema(t, db, "schema_b")
	CreateSchema(t, db, "schema_c")
	defer DropSchemas(t, db, "schema_a", "schema_b", "schema_c")

	schemas, err := db.ParseSchemas()
	require.NoError(t, err, "Failed to parse schemas")
	require.NotEmpty(t, schemas, "No schemas returned")

	schemaNames := make(map[string]bool)
	for _, schema := range schemas {
		schemaNames[schema.Name] = true

		assert.NotNil(t, schema.Database, "Schema %s has nil Database reference", schema.Name)
		assert.Equal(t, db, schema.Database, "Schema %s references wrong database", schema.Name)
	}

	expectedSchemas := []string{"schema_a", "schema_b", "schema_c"}
	for _, expected := range expectedSchemas {
		assert.True(t, schemaNames[expected], "Expected schema %s not found", expected)
	}

	t.Logf("Successfully parsed %d schemas", len(schemas))
}

func TestParseSchemasEmpty(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	schemas, err := db.ParseSchemas()
	require.NoError(t, err, "Failed to parse schemas")

	assert.NotEmpty(t, schemas, "Expected at least system schemas")

	foundPublic := false
	for _, schema := range schemas {
		if schema.Name == "public" {
			foundPublic = true
			break
		}
	}
	assert.True(t, foundPublic, "Expected 'public' schema not found")
}

func TestSchemaGetDatabase(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "test_schema")
	defer DropSchemas(t, db, "test_schema")

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	var testSchema *Schema
	for _, schema := range schemas {
		if schema.Name == "test_schema" {
			testSchema = schema
			break
		}
	}

	require.NotNil(t, testSchema, "test_schema not found")

	retrievedDB := testSchema.GetDatabase()
	assert.NotNil(t, retrievedDB, "GetDatabase returned nil")
	assert.Equal(t, db, retrievedDB, "GetDatabase returned wrong database")
}

func TestSchemasFindTable(t *testing.T) {
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

	tables, err := testSchema.LoadTables()
	require.NoError(t, err)
	require.NotEmpty(t, tables)

	userTable := testSchema.FindTable("users")
	assert.NotNil(t, userTable, "FindTable should find 'users' table")
	if userTable != nil {
		assert.Equal(t, "users", userTable.Name)
		assert.Equal(t, testSchema, userTable.Schema)
	}

	productTable := testSchema.FindTable("products")
	assert.NotNil(t, productTable, "FindTable should find 'products' table")
	if productTable != nil {
		assert.Equal(t, "products", productTable.Name)
	}

	nonExistent := testSchema.FindTable("non_existent")
	assert.Nil(t, nonExistent, "FindTable should return nil for non-existent table")
}

func TestSchemasOrder(t *testing.T) {
	db, cleanup := SetupTestDatabase(t)
	defer cleanup()

	CreateSchema(t, db, "zzz_schema")
	CreateSchema(t, db, "aaa_schema")
	defer DropSchemas(t, db, "zzz_schema", "aaa_schema")

	schemas, err := db.ParseSchemas()
	require.NoError(t, err)

	publicIndex := -1
	for i, schema := range schemas {
		if schema.Name == "public" {
			publicIndex = i
			break
		}
	}

	assert.Equal(t, 0, publicIndex, "Public schema should be sorted first")
}
