package web

import (
	"encoding/json"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GeneralResponse represents the common response body for JSON type.
type GeneralResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error response body for JSON type.
type ErrorResponse struct {
	Error  string                 `json:"error"`
	Fields map[string]interface{} `json:"fields,omitempty" swaggertype:"object"`
}

func WriteErrResponse(w http.ResponseWriter, log *zap.SugaredLogger, err error) error {
	switch errors.Cause(err).(type) {
	case validation.Errors, validation.Error:
		if errors.Is(err, bareknews.ErrDataNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			err = errors.Cause(err)
		}
	default:
		log.Error(err.Error())
		err = bareknews.ErrInternalServer
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(
		ErrorResponse{
			Error: err.Error(),
		},
	)
}

// Respond sends a JSON response to the client
func Respond(w http.ResponseWriter, data any, statusCode int) error {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}
