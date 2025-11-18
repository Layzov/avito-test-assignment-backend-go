package response

import (
	"fmt"
	"strings"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"resp_status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

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

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is required", err.Field()))
		case "min":
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' must be at least %s characters long", err.Field(), err.Param()))
		case "max":
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' must be at most %s characters long", err.Field(), err.Param()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
