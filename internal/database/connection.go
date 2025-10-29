package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type DatabaseConnection struct {
	Host       string    `json:"host"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Port       int       `json:"port"`
	Database   string    `json:"database"`
	Connected  bool      `json:"-"`
	connection *pgx.Conn `json:"-"`
	Schemas    []*Schema `json:"-"`
}

var Connections []*DatabaseConnection

func NewDatabaseConnection(host, username, password string, port int, database string) *DatabaseConnection {
	conn := &DatabaseConnection{
		Host:      host,
		Username:  username,
		Password:  password,
		Port:      port,
		Database:  database,
		Connected: false,
		// Connection: nil,
	}
	Connections = append(Connections, conn)
	return conn
}

func (dc *DatabaseConnection) Connect() error {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dc.Username, dc.Password, dc.Host, dc.Port, dc.Database)
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return err
	}
	dc.connection = conn
	dc.Connected = true
	return nil
}

func (dc *DatabaseConnection) Disconnect() error {
	if dc.Connected {
		dc.Connected = false
		return dc.connection.Close(context.Background())
	}
	return nil
}

func (dc *DatabaseConnection) findIndex(existingConnections []*DatabaseConnection) int {
	for i, conn := range existingConnections {
		if conn.Host == dc.Host && conn.Username == dc.Username && conn.Port == dc.Port && conn.Database == dc.Database {
			return i
		}
	}
	return -1
}

func (dc *DatabaseConnection) SaveAndConnect() error {
	var existingConnections []*DatabaseConnection
	if dc.Connected {
		dc.Disconnect()
	}
	err := dc.Connect()
	if err != nil {
		return err
	}
	file, err := os.ReadFile(".connections.json")
	if err != nil {
		existingConnections = []*DatabaseConnection{}
	} else {
		err = json.Unmarshal(file, &existingConnections)
		if err != nil {
			fmt.Println("Error reading connections file:", err)
		}
	}
	if idx := dc.findIndex(existingConnections); idx != -1 {
		existingConnections[idx] = dc
	} else {
		existingConnections = append(existingConnections, dc)
	}
	updatedData, err := json.MarshalIndent(existingConnections, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling connections:", err)
	}
	err = os.WriteFile(".connections.json", updatedData, 0o644)
	if err != nil {
		fmt.Println("Error writing connections file:", err)
	}
	return nil
}

func (dc *DatabaseConnection) Test() (bool, string) {
	err := dc.Connect()
	if err != nil {
		return false, err.Error()
	}
	defer dc.Disconnect()
	version := ""
	err = dc.connection.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		version = "Get version failed: " + err.Error()
		return false, version
	}

	return true, version
}

func LoadConnections() error {
	file, err := os.ReadFile(".connections.json")
	if err != nil {
		return nil
	}

	err = json.Unmarshal(file, &Connections)
	if err != nil {
		return err
	}

	return nil
}
