package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (db *Database) Connect() error {
	if db.Connected {
		return nil
	}
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.Username, db.Password, db.Host, db.Port, db.Database)
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return err
	}
	db.Connection = conn
	db.Connected = true
	return nil
}

func (db *Database) Disconnect() error {
	if db.Connected {
		db.Connected = false
		return db.Connection.Close(context.Background())
	}
	return nil
}

// SaveAndConnect connects to the database and saves it to the registry
func (db *Database) SaveAndConnect(registry *DBRegistry, configPath string) error {
	if db.Connected {
		db.Disconnect()
	}
	err := db.Connect()
	if err != nil {
		return err
	}

	existing := registry.Find(db.Host, db.Database, db.Username, db.Port)
	if existing == nil {
		registry.Add(db)
	} else {
		existing.Password = db.Password
		existing.Connected = db.Connected
		existing.Connection = db.Connection
	}

	return registry.SaveToFile(configPath)
}

func (db *Database) Test() (bool, string) {
	err := db.Connect()
	if err != nil {
		return false, err.Error()
	}
	defer db.Disconnect()
	version := ""
	err = db.Connection.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		version = "Get version failed: " + err.Error()
		return false, version
	}

	return true, version
}
