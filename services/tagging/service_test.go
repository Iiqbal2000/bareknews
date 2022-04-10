package tagging_test

import (
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/domain/tags"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tagging.New(store)
		_, err := svc.Create("tag 1")

		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.SaveCalls()), 1)
	})

	t.Run("invalid payload: tag name is blank", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tagging.New(store)
		_, err := svc.Create("")

		is := is.New(t)
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
	})

	t.Run("invalid payload: tag name is too long", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tagging.New(store)
		_, err := svc.Create("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")
		is := is.New(t)
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
	})

}

func TestUpdate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		tg := tags.New("tag 1")

		store := &tags.RepositoryMock{
			GetByIdFunc: func(id uuid.UUID) (*tags.Tags, error) {
				return tg, nil
			},
			UpdateFunc: func(tagsIn tags.Tags) error {
				tg = &tagsIn
				return nil
			},
		}

		svc := tagging.New(store)
		name := "tag 2"
		got, err := svc.Update(tg.Label.ID, name)
		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.UpdateCalls()), 1)
		is.Equal(got.Name, name)
	})

	t.Run("invalid payload: the tags is not found", func(t *testing.T) {
		store := &tags.RepositoryMock{
			GetByIdFunc: func(id uuid.UUID) (*tags.Tags, error) {
				return nil, sql.ErrNoRows
			},
			UpdateFunc: func(tagsIn tags.Tags) error {
				return nil
			},
		}

		svc := tagging.New(store)
		_, err := svc.Update(uuid.New(), "tag 2")
		is := is.New(t)
		is.Equal(err, bareknews.ErrDataNotFound)
		is.Equal(len(store.UpdateCalls()), 0)
	})
}

func TestDelete(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		store := &tags.RepositoryMock{
			CountFunc: func(id uuid.UUID) (int, error) {
				return 1, nil
			},
			DeleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}

		svc := tagging.New(store)
		is := is.New(t)

		err := svc.Delete(uuid.New())
		is.Equal(err, nil)
		is.Equal(len(store.DeleteCalls()), 1)
	})

	t.Run("invalid payload: the tags is not found", func(t *testing.T) {
		store := &tags.RepositoryMock{
			CountFunc: func(id uuid.UUID) (int, error) {
				return 0, sql.ErrNoRows
			},
			DeleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}

		svc := tagging.New(store)
		is := is.New(t)

		err := svc.Delete(uuid.New())
		is.Equal(err, sql.ErrNoRows)
		is.Equal(len(store.DeleteCalls()), 0)
	})
}
