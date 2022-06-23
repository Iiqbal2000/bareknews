package web

import (
	"encoding/json"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// RespBody represents the common response body for JSON type.
type RespBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrRespBody represents an error response body for JSON type.
type ErrRespBody struct {
	Err map[string]interface{} `json:"error" swaggertype:"object"`
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
		ErrRespBody{
			Err: map[string]interface{}{
				"message": err.Error(),
			},
		},
	)
}
