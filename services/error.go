package services

import validation "github.com/go-ozzo/ozzo-validation/v4"

var ErrInternalServer = validation.NewError(
	"internal_server_error", 
	"Internal server error",
)