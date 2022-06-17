package news_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreate(t *testing.T) {
	t.Run("valid payload should be success", func(t *testing.T) {
		payloadTest := []struct {
			input        news.NewsIn
			wantSaveCall int
			wantGetCall  int
		}{
			{
				input: news.NewsIn{
					Title:  "news title",
					Body:   "news body",
					Status: "draft",
					Tags:   []string{"tag1"},
				},
				wantSaveCall: 1,
				wantGetCall:  1,
			},
			{
				input: news.NewsIn{
					Title:  "news title",
					Body:   "news body",
					Status: "Draft",
					Tags:   []string{"tag1"},
				},
				wantSaveCall: 1,
				wantGetCall:  1,
			},
			{
				input: news.NewsIn{
					Title:  "news title",
					Body:   "news body",
					Status: "publish",
					Tags:   []string{},
				},
				wantSaveCall: 1,
				wantGetCall:  1,
			},
		}

		for i, pt := range payloadTest {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
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

				svc := news.CreateSvc(store, tags.CreateSvc(tgStore))
				_, err := svc.Create(context.TODO(), pt.input)
				is.NoErr(err)
				is.Equal(len(store.SaveCalls()), pt.wantSaveCall)
				is.Equal(len(tgStore.GetByNamesCalls()), pt.wantGetCall)
			})
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		payloadTest := []struct {
			name         string
			input        news.NewsIn
			wantSaveCall int
			wantGetCall  int
		}{
			{
				name: "title of news is blank",
				input: news.NewsIn{
					Title:  " ",
					Body:   "news body",
					Status: "draft",
					Tags:   []string{"tag1"},
				},
				wantSaveCall: 0,
				wantGetCall:  1,
			},
			{
				name: "status of news is invalid",
				input: news.NewsIn{
					Title:  "news title",
					Body:   "news body",
					Status: "draftler",
					Tags:   []string{"tag1"},
				},
				wantSaveCall: 0,
				wantGetCall:  1,
			},
		}

		for _, test := range payloadTest {
			t.Run(test.name, func(t *testing.T) {
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
				svc := news.CreateSvc(store, tags.CreateSvc(tgStore))
				_, err := svc.Create(context.TODO(), test.input)
				is.True(err != nil)
				is.Equal(len(store.SaveCalls()), test.wantSaveCall)
				is.Equal(len(tgStore.GetByNamesCalls()), test.wantGetCall)
			})
		}

	})
}

func TestUpdate(t *testing.T) {
	tgId := uuid.New()
	payload := news.Create("news title", "news body", "draft", []uuid.UUID{tgId})

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
			return []tags.Tags{{Label: bareknews.Label{ID: tgId}}}, nil
		},
	}

	is := is.New(t)

	newPayload := news.NewsIn{
		Title: "news title update",
	}

	svc := news.CreateSvc(store, tags.CreateSvc(tgStore))
	resp, err := svc.Update(context.TODO(), payload.Post.ID, newPayload)
	is.NoErr(err)
	is.True(resp.Title != "news title")
	is.True(resp.Body != "")
	is.True(resp.Slug != "")
	is.True(resp.Status != "")
	is.Equal(len(resp.Tags), 1)
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

		svc := news.CreateSvc(store, tags.CreateSvc(&tags.RepositoryMock{}))
		payload := news.Create("news title", "news body", "draft", []uuid.UUID{uuid.New()})
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

		svc := news.CreateSvc(store, tags.CreateSvc(&tags.RepositoryMock{}))
		err := svc.Delete(context.TODO(), uuid.New())
		is.True(err != nil)
		is.Equal(len(store.CountCalls()), 1)
		is.Equal(len(store.DeleteCalls()), 0)
	})
}

func TestGetAllByStatus(t *testing.T) {
	nwsStore := &news.RepositoryMock{
		GetAllByStatusFunc: func(ctx context.Context, status bareknews.Status) ([]news.News, error) {
			return nil, nil
		},
	}
	tgStore := &tags.RepositoryMock{}

	nwsSvc := news.CreateSvc(nwsStore, tags.CreateSvc(tgStore))

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
