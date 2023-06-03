package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/devcashflow/database-server/pkg/api"
	"github.com/devcashflow/database-server/pkg/database/mysql"
	"github.com/devcashflow/database-server/types"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHandleCreateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mdb := &mysql.DB{db}
	server, err := handlers.New(mdb)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectPrepare("INSERT INTO emails").
		ExpectExec().
		WithArgs("test@example.com", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := &types.CreateEmailRequest{
		Email: "test@example.com",
		Name:  "test",
	}
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/create-email", bytes.NewReader(body))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleCreateEmail)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var response types.SucessResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if !response.Success {
		t.Error("handler returned unexpected body: got success=false want success=true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleListEmails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mdb := &mysql.DB{db}
	server, err := handlers.New(mdb)

	rows := sqlmock.NewRows([]string{"email"}).
		AddRow("test@example.com").
		AddRow("another_test@example.com")

	mock.ExpectPrepare("SELECT email FROM emails").
		ExpectQuery().
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/emails", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleListEmails)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, http.StatusOK)
	}

	var response types.ListEmailsResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response.Emails) != 2 {
		t.Errorf("handler returned unexpected number of emails: got %d, want 2",
			len(response.Emails))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
