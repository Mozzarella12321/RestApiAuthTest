package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type LoginResponse struct {
	Status string    `json:"status"`
	Token  uuid.UUID `json:"token"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func LoggedInOK(token uuid.UUID) LoginResponse {
	return LoginResponse{
		Status: StatusOK,
		Token:  token,
	}
}
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// Need to fix it later
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s characters long", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
