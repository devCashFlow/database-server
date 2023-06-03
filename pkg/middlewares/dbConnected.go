package middleware

import (
	"net/http"

	"github.com/devcashflow/database-server/pkg/database/mysql"
)

func DBConnected(db *mysql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := db.Ping()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
