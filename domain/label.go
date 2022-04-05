package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Label struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (l Label) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required, validation.Length(1, 20)),
	)
}
