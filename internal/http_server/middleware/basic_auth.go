package my_middleware

import (
	"net/http"
	"restapiauthtest/internal/auth"
	"restapiauthtest/lib/api/response"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Storage interface {
	SaveNewUser(login string, hash []byte) error
	GetUserData(login string) ([]byte, error)
	ExecuteQuery(query string, args ...interface{}) error
	CreateToken(token uuid.UUID) error
}

// eto kakoeto ahaha
func MyBasicAuth(s Storage, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		_, err := auth.Login(s, username, password)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	}
}
