package database

import (
	"encoding/json"
	"os"
	"sync"
)

// DBRegistry manages a collection of database connections
type DBRegistry struct {
	databases []*Database
	mu        sync.RWMutex
}

// NewDBRegistry creates a new database registry
func NewDBRegistry() *DBRegistry {
	return &DBRegistry{
		databases: make([]*Database, 0),
	}
}

// Add adds a database connection to the registry
func (r *DBRegistry) Add(db *Database) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.databases = append(r.databases, db)
}

// GetAll returns all database connections
func (r *DBRegistry) GetAll() []*Database {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.databases
}

// Find finds a database connection by host, database name, username and port
func (r *DBRegistry) Find(host, database, username string, port int) *Database {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, conn := range r.databases {
		if conn.Host == host && conn.Database == database && conn.Username == username && conn.Port == port {
			return conn
		}
	}
	return nil
}

// LoadFromFile loads database connections from a JSON file
func (r *DBRegistry) LoadFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	err = json.Unmarshal(file, &r.databases)
	if err != nil {
		return err
	}

	return nil
}

// SaveToFile saves all database connections to a JSON file
func (r *DBRegistry) SaveToFile(path string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := json.MarshalIndent(r.databases, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// Remove removes a database connection from the registry
func (r *DBRegistry) Remove(db *Database) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, conn := range r.databases {
		if conn == db {
			r.databases = append(r.databases[:i], r.databases[i+1:]...)
			return
		}
	}
}

// Count returns the number of database connections
func (r *DBRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.databases)
}
