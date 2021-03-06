package bareknews

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

// Label is an entity that represents the label
type Label struct {
	ID   uuid.UUID
	Name string
}

func (l Label) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required, validation.Length(1, 20)),
	)
}
