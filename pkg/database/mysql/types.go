package mysql

import (
	"database/sql"
	"time"
)

type DB struct {
	*sql.DB
}

// Email is the struct used to insert at the database.
type Email struct {
	id      int       `db:"id"`
	name    string    `db:"name"`
	address string    `db:"email"`
	created time.Time `db:"creeated_at"`
}
