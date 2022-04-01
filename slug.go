package bareknews

import (
	"fmt"
	"strings"
)

// Value object
type Slug string

func NewSlug(input string) Slug {
	slug := new(strings.Builder)

	fmt.Fprintf(slug, "%s", strings.ToLower(strings.Replace(input, " ", "-", -1)))

	return Slug(slug.String())
}
