package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/devcashflow/database-server/pkg/database/mysql"
	"github.com/devcashflow/database-server/pkg/database/postgresql"
	"github.com/devcashflow/database-server/types"
)

type DB struct {
	SQLDB *sql.DB
}

type Config struct {
	// the type of the db - mysql, postgresql, etc
	Type string
	//
	User string
	Pass string
	Host string
	Port string
	// Names represents an array of databases name to connect to.
	Names []string
	// Name represents databases name to connect to.
	Name string

	SSLMODE string
}

type Database interface {
	InsertEmail(email *types.Email) error
	ListEmails() ([]types.Email, error)
	Version() (types.Version, error)
	Ping() error
	// ... add more methods as needed
}

func Connect(config Config) (Database, error) {
	if config.User == "" {
		return nil, errors.New("Database User is required")
	}

	if config.Pass == "" {
		return nil, errors.New("Database Pass is required")
	}

	if config.Host == "" {
		return nil, errors.New("Database Host is required")
	}

	if config.Name == "" {
		return nil, errors.New("Database Name is required")
	}

	if config.Type == "" {
		return nil, errors.New("Database Type is required")
	}
	user := config.User
	pass := config.Pass
	host := config.Host
	port := config.Port
	dbName := config.Name
	dbType := config.Type
	sslmode := config.SSLMODE

	var db *sql.DB
	var err error

	switch dbType {
	case "mysql":
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbName)
		db, err = sql.Open("mysql", connStr)
	case "postgres":
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s door=%s sslmode=%s", user, pass, dbName, host, port, sslmode)
		db, err = sql.Open("postgres", connStr)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	switch dbType {
	case "mysql":
		return &mysql.MySQLDB{SQLDB: db}, nil
	case "postgres":
		return &postgresql.PostgreSQLDB{SQLDB: db}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
