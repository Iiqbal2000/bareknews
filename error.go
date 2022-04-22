package bareknews

import (
	"encoding/json"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

var ErrInternalServer = errors.New("internal server error")

var ErrDataAlreadyExist = validation.NewError(
	"data_already_exist",
	"the data already exist",
)

var ErrDataNotFound = validation.NewError(
	"data_not_found",
	"the data not found",
)

var ErrInvalidJSON = validation.NewError(
	"invalid_json",
	"the syntax JSON is invalid",
)

const SubStrUniqueConstraint = "UNIQUE constraint failed:"

func WriteErrResponse(w http.ResponseWriter, err error) error {
	switch errors.Cause(err).(type) {
	case validation.Errors, validation.Error:
		if err.Error() == ErrDataNotFound.Error() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		log.Println(err.Error())
		err = ErrInternalServer
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
