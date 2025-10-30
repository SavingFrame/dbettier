package database

import (
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
}

func NewDatabase(host, username, password string, port int, database string) *Database {
	return &Database{
		Host:      host,
		Username:  username,
		Password:  password,
		Port:      port,
		Database:  database,
		Connected: false,
	}
}
