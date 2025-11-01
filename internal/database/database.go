package database

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Host       string    `json:"host"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Port       int       `json:"port"`
	Database   string    `json:"database"`
	Connected  bool      `json:"-"`
	connection *pgx.Conn `json:"-"`
	Schemas    []*Schema `json:"-"`
	ID         string    `json:"id"`
}

func NewDatabase(host, username, password string, port int, database string) *Database {
	return &Database{
		Host:      host,
		Username:  username,
		Password:  password,
		Port:      port,
		Database:  database,
		Connected: false,
		ID:        generateDatabaseID(host, username, port, database),
	}
}

func generateDatabaseID(host, username string, port int, database string) string {
	input := fmt.Sprintf("%s:%s:%d:%s", host, username, port, database)

	hash := sha256.Sum256([]byte(input))

	return hex.EncodeToString(hash[:])
}
