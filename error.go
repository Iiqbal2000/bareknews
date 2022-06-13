package bareknews

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrDataAlreadyExist = validation.NewError(
		"data_already_exist",
		"the data already exist",
	)
	ErrDataNotFound = validation.NewError(
		"data_not_found",
		"the data is not found",
	)
	ErrInvalidJSON = validation.NewError(
		"invalid_json",
		"the JSON syntax is invalid",
	)
)

const SubStrUniqueConstraint = "UNIQUE constraint failed:"
