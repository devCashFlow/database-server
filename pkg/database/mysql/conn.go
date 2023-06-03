package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(user, pass, host, door, dbName string) (*DB, error) {
	db, err := dbConn(user, pass, host, door, dbName)
	if err != nil {
		return nil, err
	}
	// Attempt to ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if the ping fails
		return nil, err
	}
	database := &DB{db}

	return database, nil
}

func dbConn(user, pass, host, door, dbName string) (*sql.DB, error) {
	// fmt.Printf("user: %v, pass: %v, host: %v, door: %v, dbName: %v", user, pass, host, door, dbName)
	db, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+":"+door+")/"+dbName+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	return db, nil
}
