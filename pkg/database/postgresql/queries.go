package postgresql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/devcashflow/database-server/types"
)

type PostgreSQLDB struct {
	SQLDB *sql.DB
}

func (db *PostgreSQLDB) Version() (types.Version, error) {
	rows, err := db.SQLDB.Query("select version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}
	return types.Version{
		Version: version,
	}, nil
}

// InsertEmail inserts an email at the database. It assumes sanity checks are
// done before getting here, so it only inserts.
func (db *PostgreSQLDB) InsertEmail(email *types.Email) error {
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

func (db *PostgreSQLDB) ListEmails() ([]types.Email, error) {
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

func (db *PostgreSQLDB) Ping() error {
	return db.SQLDB.Ping()
}
