package registration

import (
	"log/slog"
	"net/http"
	"restapiauthtest/internal/auth"
	"restapiauthtest/lib/api/response"
	"restapiauthtest/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Storage interface {
	SaveNewUser(login string, hash []byte) error
	GetUserData(login string) ([]byte, error)
	ExecuteQuery(query string, args ...interface{}) error
	CreateToken(token uuid.UUID) error
}

type User struct {
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}

func Registration(log *slog.Logger, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Registration"

		log = log.With(
			slog.String("op", op),
			slog.String("ReqID", middleware.GetReqID(r.Context())),
		)

		var user User
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		log.Info("request body decoded", slog.Any("request", user))

		if err := validator.New().Struct(user); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		err := auth.RegisterNewUser(s, user.Login, user.Password)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.OK())
	}
}
