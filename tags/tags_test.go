package tags_test

import (
	"testing"

	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/matryer/is"
)

func TestNewTags(t *testing.T) {
	t.Run("Valid tags", func(t *testing.T) {
		tag := tags.Create("tag 1")
		err := tag.Validate()
		is := is.New(t)
		is.NoErr(err)
		is.True(tag.Label.ID.String() != "")
		is.True(tag.Label.Name != "")
		is.True(tag.Slug != "")
	})

	t.Run("Invalid tags", func(t *testing.T) {
		tag := tags.Create("")
		err := tag.Validate()
		is := is.New(t)
		is.True(err != nil)
	})
}

func TestChangeNameTags(t *testing.T) {
	t.Run("Valid change", func(t *testing.T) {
		is := is.New(t)
		tag := tags.Create("tag1")
		err := tag.Validate()
		is.NoErr(err)

		oldTag := *tag

		tag.ChangeName("tag2")
		err = tag.Validate()
		is.NoErr(err)
		is.True(tag.Label.Name != oldTag.Label.Name)
		is.True(tag.Slug != oldTag.Slug)
	})

	t.Run("Invalid change", func(t *testing.T) {
		is := is.New(t)
		tag := tags.Create("tag1")
		err := tag.Validate()
		is.NoErr(err)

		tag.ChangeName("")
		err = tag.Validate()
		is.True(err != nil)
	})
}
