package mysql_test

import (
	"testing"

	"github.com/devcashflow/database-server/pkg/database/mysql"
	"github.com/devcashflow/database-server/types"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mdb := &mysql.DB{db}

	mock.ExpectPrepare("INSERT INTO emails").
		ExpectExec().
		WithArgs("test@example.com", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = mdb.InsertEmail(types.Email{Address: "test@example.com", Name: "test"})
	if err != nil {
		t.Errorf("error was not expected while inserting email: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mdb := &mysql.DB{db}

	rows := sqlmock.NewRows([]string{"email"}).
		AddRow("test@example.com").
		AddRow("another_test@example.com")

	mock.ExpectPrepare("SELECT email FROM emails").
		ExpectQuery().
		WillReturnRows(rows)

	emails, err := mdb.ListEmails()
	if err != nil {
		t.Errorf("error was not expected while listing emails: %s", err)
	}

	if len(emails) != 2 {
		t.Errorf("expected 2 emails, got %d", len(emails))
	}

	if emails[0].Address != "test@example.com" {
		t.Errorf("expected first email to be 'test@example.com', got '%s'", emails[0].Address)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
