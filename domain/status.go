package domain

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Status is a value object that represents the status of the news.
type Status string

const (
	Publish Status = "publish"
	Draft   Status = "draft"
	Deleted Status = "deleted"
)

// Validate performs validating to the status.
func (s Status) Validate() error {
	return validation.Validate(
		s.String(),
		validation.Required.Error("status cannot be blank"),
		validation.In(
			Publish.String(), 
			Draft.String(), 
			Deleted.String(),
			).Error("status must be one of 'publish', 'draft', 'deleted'"),
	)
}

func (s Status) String() string {
	return string(s)
}
