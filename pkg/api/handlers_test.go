package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/devcashflow/database-server/pkg/api"
	"github.com/devcashflow/database-server/types"
)

type MockDB struct {
	InsertEmailError error
	ListEmailsValue  []types.Email
	ListEmailsError  error
	VersionValue     types.Version
	VersionError     error
}

func (mdb *MockDB) Close() error {
	// Implement the Close method of the Database interface for the mock
	return nil
}
func (m *MockDB) InsertEmail(email *types.Email) error {
	return nil
}

func (m *MockDB) ListEmails() ([]types.Email, error) {
	// Return the values specified in the mock
	return m.ListEmailsValue, m.ListEmailsError
}

func (m *MockDB) Version() (types.Version, error) {
	// Return the values specified in the mock
	return m.VersionValue, m.VersionError
}

func (m *MockDB) Ping() error {
	return nil
}

func TestHandleCreateEmail(t *testing.T) {
	mdb := &MockDB{}

	server, err := handlers.New(mdb)
	if err != nil {
		t.Fatalf("unexpected error creating server: %s", err)
	}

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
}

func TestHandleListEmails(t *testing.T) {
	mdb := &MockDB{
		ListEmailsValue: []types.Email{
			{Address: "test@example.com"},
			{Address: "another_test@example.com"},
		},
	}

	server, err := handlers.New(mdb)
	if err != nil {
		t.Fatalf("unexpected error creating server: %s", err)
	}

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
}
