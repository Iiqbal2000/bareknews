package posting_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/Iiqbal2000/bareknews/domain/tags"
	"github.com/Iiqbal2000/bareknews/services/posting"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		store := &news.RepositoryMock{
			SaveFunc: func(ctx context.Context, news news.News) error {
				return nil
			},
		}

		tgStore := &tags.RepositoryMock{
			GetByNamesFunc: func(ctx context.Context, names ...string) ([]tags.Tags, error) {
				return nil, nil
			},
		}

		is := is.New(t)

		svc := posting.New(store, tagging.New(tgStore))
		_, err := svc.Create(context.TODO(), "news title", "news body", "draft", []string{"tag1"})
		is.NoErr(err)
		is.Equal(len(store.SaveCalls()), 1)
		is.Equal(len(tgStore.GetByNamesCalls()), 1)
	})

	t.Run("invalid payload: news title is blank", func(t *testing.T) {
		store := &news.RepositoryMock{
			SaveFunc: func(ctx context.Context, news news.News) error {
				return nil
			},
		}

		tgStore := &tags.RepositoryMock{
			GetByNamesFunc: func(ctx context.Context, names ...string) ([]tags.Tags, error) {
				return nil, nil
			},
		}

		is := is.New(t)

		svc := posting.New(store, tagging.New(tgStore))
		_, err := svc.Create(context.TODO(), "", "news body", "draft", []string{"tag1"})
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
		is.Equal(len(tgStore.GetByNamesCalls()), 1)
	})

	t.Run("invalid payload: news status is wrong", func(t *testing.T) {
		store := &news.RepositoryMock{
			SaveFunc: func(ctx context.Context, news news.News) error {
				return nil
			},
		}

		tgStore := &tags.RepositoryMock{
			GetByNamesFunc: func(ctx context.Context, names ...string) ([]tags.Tags, error) {
				return nil, nil
			},
		}

		is := is.New(t)

		svc := posting.New(store, tagging.New(tgStore))
		_, err := svc.Create(context.TODO(), "news title", "news body", "draftler", []string{"tag1"})
		is.True(err != nil)
		is.Equal(len(store.SaveCalls()), 0)
		is.Equal(len(tgStore.GetByNamesCalls()), 1)
	})
}

func TestUpdate(t *testing.T) {
	payload := news.New("news title", "news body", "draft", []uuid.UUID{uuid.New()})

	store := &news.RepositoryMock{
		GetByIdFunc: func(ctx context.Context, id uuid.UUID) (*news.News, error) {
			return payload, nil
		},
		UpdateFunc: func(ctx context.Context, news news.News) error {
			return nil
		},
	}

	tgStore := &tags.RepositoryMock{
		GetByNamesFunc: func(ctx context.Context, names ...string) ([]tags.Tags, error) {
			return nil, nil
		},
		GetByIdsFunc: func(ctx context.Context, ids []uuid.UUID) ([]tags.Tags, error) {
			return nil, nil
		},
	}

	is := is.New(t)

	svc := posting.New(store, tagging.New(tgStore))
	_, err := svc.Update(context.TODO(), payload.Post.ID, "news title update", "", "", []string{})
	is.NoErr(err)
	is.Equal(len(store.GetByIdCalls()), 1)
	is.Equal(len(store.UpdateCalls()), 1)
	is.Equal(len(tgStore.GetByNamesCalls()), 0)
}

func TestDelete(t *testing.T) {
	t.Run("valid input should be success", func(t *testing.T) {
		store := &news.RepositoryMock{
			CountFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
				return 1, nil
			},
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}

		is := is.New(t)

		svc := posting.New(store, tagging.New(&tags.RepositoryMock{}))
		payload := news.New("news title", "news body", "draft", []uuid.UUID{uuid.New()})
		err := svc.Delete(context.TODO(), payload.Post.ID)
		is.NoErr(err)

		is.Equal(len(store.CountCalls()), 1)
		is.Equal(len(store.DeleteCalls()), 1)
	})

	t.Run("invalid payload: the news is not found", func(t *testing.T) {
		store := &news.RepositoryMock{
			CountFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
				return 0, sql.ErrNoRows
			},
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}

		is := is.New(t)

		svc := posting.New(store, tagging.New(&tags.RepositoryMock{}))
		err := svc.Delete(context.TODO(), uuid.New())
		is.True(err != nil)
		is.Equal(err.Error(), sql.ErrNoRows.Error())
		is.Equal(len(store.CountCalls()), 1)
		is.Equal(len(store.DeleteCalls()), 0)
	})
}

func TestGetAllByStatus(t *testing.T) {
	nwsStore := &news.RepositoryMock{
		GetAllByStatusFunc: func(ctx context.Context, status domain.Status) ([]news.News, error) {
			return nil, nil
		},
	}
	tgStore := &tags.RepositoryMock{}

	nwsSvc := posting.New(nwsStore, tagging.New(tgStore))

	is := is.New(t)
	_, err := nwsSvc.GetAllByStatus(context.TODO(), "draft")
	is.NoErr(err)

	_, err = nwsSvc.GetAllByStatus(context.TODO(), "publish")
	is.NoErr(err)

	_, err = nwsSvc.GetAllByStatus(context.TODO(), "")
	is.True(err != nil)

	_, err = nwsSvc.GetAllByStatus(context.TODO(), "publsjsja")
	is.True(err != nil)
}
