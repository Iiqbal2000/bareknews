package tagging

import (
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/domain/tag"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		store := &tag.RepositoryMock{
			SaveFunc: func(tags bareknews.Tags) error {
				return nil
			},
		}

		svc := New(store)
		_, err := svc.Create("teg 1")

		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.SaveCalls()), 1)
	})

	t.Run("Invalid: blank tag name", func(t *testing.T) {
		store := &tag.RepositoryMock{
			SaveFunc: func(tags bareknews.Tags) error {
				return nil
			},
		}

		svc := New(store)
		_, err := svc.Create("")

		is := is.New(t)
		is.Equal(err, bareknews.ErrBlankTag)
		is.Equal(len(store.SaveCalls()), 0)
	})

	t.Run("Invalid: tag name too long", func(t *testing.T) {
		store := &tag.RepositoryMock{
			SaveFunc: func(tags bareknews.Tags) error {
				return nil
			},
		}

		svc := New(store)
		_, err := svc.Create("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")
		is := is.New(t)
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		tg, _ := bareknews.NewTags("tag 1")

		store := &tag.RepositoryMock{
			GetByIdFunc: func(id string) (*bareknews.Tags, error) {
				return tg, nil
			},
			UpdateFunc: func(tagsIn bareknews.Tags) error {
				tg = &tagsIn
				return nil
			},
		}

		svc := New(store)
		err := svc.Update(tg.ID.String(), "tag 2")
		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.UpdateCalls()), 1)
	})

	t.Run("Invalid: Not found", func(t *testing.T) {
		store := &tag.RepositoryMock{
			GetByIdFunc: func(id string) (*bareknews.Tags, error) {
				return nil, sql.ErrNoRows
			},
			UpdateFunc: func(tagsIn bareknews.Tags) error {
				return nil
			},
		}

		svc := New(store)
		err := svc.Update("tag item is not found", "tag 2")
		is := is.New(t)
		is.Equal(err, sql.ErrNoRows)
		is.Equal(len(store.UpdateCalls()), 0)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		store := &tag.RepositoryMock{
			GetByIdFunc: func(s string) (*bareknews.Tags, error) {
				return nil, nil
			},
			DeleteFunc: func(s string) error {
				return nil
			},
		}

		svc := New(store)
		is := is.New(t)

		err := svc.Delete(uuid.New().String())
		is.Equal(err, nil)
		is.Equal(len(store.DeleteCalls()), 1)
	})
}
