package restapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func WriteErrResponse(w http.ResponseWriter, err error) error {
	switch errors.Cause(err).(type) {
	case validation.Errors, validation.Error:
		if err.Error() == bareknews.ErrDataNotFound.Error() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		log.Println(err.Error())
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
