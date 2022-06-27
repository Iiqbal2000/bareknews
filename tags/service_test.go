package tags_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(ctx context.Context, tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tags.CreateSvc(store)
		_, err := svc.Create(context.TODO(), "tag 1")

		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.SaveCalls()), 1)
	})

	t.Run("invalid payload: tag name is blank", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(ctx context.Context, tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tags.CreateSvc(store)
		_, err := svc.Create(context.TODO(), "")

		is := is.New(t)
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
	})

	t.Run("invalid payload: tag name is too long", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(ctx context.Context, tags tags.Tags) error {
				return nil
			},
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tags.CreateSvc(store)
		_, err := svc.Create(context.TODO(), "Lorem Ipsum is simply dummy text of the printing and typesetting industry.")
		is := is.New(t)
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
	})

	t.Run("invalid payload: tag name already exists", func(t *testing.T) {
		store := &tags.RepositoryMock{
			SaveFunc: func(ctx context.Context, tags tags.Tags) error {
				return bareknews.ErrDataAlreadyExist
			},
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return &tags.Tags{}, nil
			},
		}

		svc := tags.CreateSvc(store)
		_, err := svc.Create(context.TODO(), "Lorem Ipsum")
		is := is.New(t)
		is.True(err != nil)
		is.Equal(err, bareknews.ErrDataAlreadyExist)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		tg := tags.Create("tag 1")

		store := &tags.RepositoryMock{
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return tg, nil
			},
			UpdateFunc: func(ctx context.Context, tagsIn tags.Tags) error {
				tg = &tagsIn
				return nil
			},
		}

		svc := tags.CreateSvc(store)
		name := "tag 2"
		got, err := svc.Update(context.TODO(), tg.Label.ID, name)
		is := is.New(t)
		is.Equal(err, nil)
		is.Equal(len(store.UpdateCalls()), 1)
		is.Equal(got.Name, name)
	})

	t.Run("invalid payload: the tags is not found", func(t *testing.T) {
		store := &tags.RepositoryMock{
			GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
				return nil, sql.ErrNoRows
			},
			UpdateFunc: func(ctx context.Context, tagsIn tags.Tags) error {
				return nil
			},
		}

		svc := tags.CreateSvc(store)
		_, err := svc.Update(context.TODO(), uuid.New(), "tag 2")
		is := is.New(t)
		is.Equal(err, bareknews.ErrDataNotFound)
		is.Equal(len(store.UpdateCalls()), 0)
	})
}

func TestDelete(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		store := &tags.RepositoryMock{
			CountFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
				return 1, nil
			},
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}

		svc := tags.CreateSvc(store)
		is := is.New(t)

		err := svc.Delete(context.TODO(), uuid.New())
		is.Equal(err, nil)
		is.Equal(len(store.DeleteCalls()), 1)
	})

	t.Run("invalid payload: the tags is not found", func(t *testing.T) {
		store := &tags.RepositoryMock{
			CountFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
				return 0, sql.ErrNoRows
			},
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}

		svc := tags.CreateSvc(store)
		is := is.New(t)

		err := svc.Delete(context.TODO(), uuid.New())
		is.True(err != nil)
		is.Equal(len(store.DeleteCalls()), 0)
	})
}
