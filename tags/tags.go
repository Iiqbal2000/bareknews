package tags

import (
	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

// Tags is an aggregrate that represent a tag
type Tags struct {
	Label bareknews.Label
	Slug  bareknews.Slug
}

func Create(tagName string) *Tags {
	tag := Tags{
		Label: bareknews.Label{
			ID:   uuid.New(),
			Name: tagName,
		},
		Slug: bareknews.NewSlug(tagName),
	}

	return &tag
}

func (t *Tags) ChangeName(newName string) {
	t.Label.Name = newName
	t.Slug = bareknews.NewSlug(newName)
}

func (t Tags) Validate() error {
	return t.Label.Validate()
}
