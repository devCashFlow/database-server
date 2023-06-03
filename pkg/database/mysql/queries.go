package mysql

import (
	"github.com/devcashflow/database-server/types"

	_ "github.com/go-sql-driver/mysql"
)

// InsertEmail inserts an email at the database. It assumes sanity checks are
// done before getting here, so it only inserts.
func (db *DB) InsertEmail(email types.Email) error {
	insForm, err := db.Prepare("INSERT INTO emails (email, name) VALUES (?,?)")
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

func (db *DB) ListEmails() ([]types.Email, error) {
	stmt, err := db.Prepare("SELECT email FROM emails")
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
		var email Email
		if err := rows.Scan(&email.address); err != nil {
			return nil, err
		}
		emails = append(emails, types.Email{
			Address: email.address,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
