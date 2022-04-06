package storage

import (
	"testing"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

const dbfile = "./../../bareknews.db"

func TestSaveNews(t *testing.T) {
	conn := Run(dbfile, true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	want := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(*want)
	is.NoErr(err)

	got, err := newsStore.GetById(want.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Status, want.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
}

func TestUpdateNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	is := is.New(t)
	tgId := uuid.New()

	news := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(*news)
	is.NoErr(err)

	wantTitle := gofakeit.Sentence(10)
	wantStatus := domain.Publish

	news.ChangeTitle(wantTitle)
	news.ChangeStatus(wantStatus)

	err = newsStore.Update(*news)
	is.NoErr(err)

	got, err := newsStore.GetById(news.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, wantTitle)
	is.Equal(got.Status, wantStatus)
}

func TestDeleteNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	news := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(*news)
	is.NoErr(err)

	err = newsStore.Delete(news.Post.ID)
	is.NoErr(err)

	_, err = newsStore.GetById(news.Post.ID)
	is.True(err != nil)
}

func TestGetByID(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	want := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(*want)
	is.NoErr(err)

	got, err := newsStore.GetById(want.Post.ID)

	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Status, want.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
	is.Equal(len(got.TagsID), 1)
}

func TestGetAllNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}

	tgId := uuid.New()

	wantNews1 := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(*wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	wantNews2 := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []uuid.UUID{tgId})
	err = newsStore.Save(*wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAll()
	is := is.New(t)
	is.NoErr(err)
	is.Equal(len(got), 2)

	is.Equal(got[0].Post.Title, wantNews1.Post.Title)
	is.Equal(got[0].Status, wantNews1.Status)
	is.Equal(got[0].Post.Body, wantNews1.Post.Body)
	is.Equal(got[0].Slug, wantNews1.Slug)
	is.Equal(len(got[0].TagsID), 1)

	is.Equal(got[1].Post.Title, wantNews2.Post.Title)
	is.Equal(got[1].Status, wantNews2.Status)
	is.Equal(got[1].Post.Body, wantNews2.Post.Body)
	is.Equal(got[1].Slug, wantNews2.Slug)
	is.Equal(len(got[1].TagsID), 1)
}
