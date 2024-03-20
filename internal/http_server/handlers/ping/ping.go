package ping

import (
	"log/slog"
	"net/http"
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
	// ExecuteQuery(query string, args ...interface{}) error
	CreateToken(token uuid.UUID) error
	GetToken(token uuid.UUID) error
}

type Request struct {
	Token uuid.UUID `json:"token" validate:"required"`
}

func Ping(log *slog.Logger, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"

		log = log.With(
			slog.String("op", op),
			slog.String("ReqID", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		if req.Token == uuid.Nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("Unauthorized"))
		}
		err := Storage.GetToken(s, req.Token)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(err.Error()))

			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Pong())
	}
}
