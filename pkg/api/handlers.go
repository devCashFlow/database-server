package handlers

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/devcashflow/database-server/pkg/database"
	"github.com/devcashflow/database-server/types"

	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

type Server struct {
	db database.Database
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func New(db database.Database) (*Server, error) {
	if db == nil {
		return nil, errors.New("No db informed")
	}
	return &Server{
		db: db,
	}, nil
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServer(err error) render.Renderer {
	// log err at log file
	fmt.Printf("err: %+v\n\n", err)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server Error",
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// CreateEmail stores a new email at the database.
func (s *Server) HandleCreateEmail(w http.ResponseWriter, r *http.Request) {
	var request types.CreateEmailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}
	if request.Email == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("Email Address Missing")))
		return
	}
	err = s.db.InsertEmail(&types.Email{
		Address: request.Email,
		Name:    request.Name,
	})
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.SucessResponse{
		Success: true,
	})
}

func (s *Server) HandleListEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := s.db.ListEmails()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.ListEmailsResponse{
		Emails: emails,
	})
}

func (s *Server) HandleVersion(w http.ResponseWriter, r *http.Request) {
	version, err := s.db.Version()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.Version{
		Version: version.Version,
	})
}
