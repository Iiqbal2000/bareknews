package bareknews

import (
	"fmt"
	"strings"
)

// Slug is a value object that represents the slug
type Slug string

func NewSlug(input string) Slug {
	slug := new(strings.Builder)

	fmt.Fprintf(slug, "%s", strings.ToLower(strings.Replace(input, " ", "-", -1)))

	return Slug(slug.String())
}

func (s Slug) String() string {
	return string(s)
}