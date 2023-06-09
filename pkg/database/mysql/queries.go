package mysql

import (
	"database/sql"
	"log"

	"github.com/devcashflow/database-server/types"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	SQLDB *sql.DB
}

// InsertEmail inserts an email at the database. It assumes sanity checks are
// done before getting here, so it only inserts.
func (db *MySQLDB) InsertEmail(email *types.Email) error {
	insForm, err := db.SQLDB.Prepare("INSERT INTO emails (email, name) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer insForm.Close()

	_, err = insForm.Exec(email.Address, email.Name)
	if err != nil {
		return err
	}

	return nil
}

func (db *MySQLDB) ListEmails() ([]types.Email, error) {
	stmt, err := db.SQLDB.Prepare("SELECT email FROM emails")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var emails []types.Email
	for rows.Next() {
		var email types.Email
		if err := rows.Scan(&email.Address); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func (db *MySQLDB) Version() (types.Version, error) {
	rows, err := db.SQLDB.Query("select version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var version types.Version
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}
	return version, nil
}

func (db *MySQLDB) Ping() error {
	return db.SQLDB.Ping()
}
