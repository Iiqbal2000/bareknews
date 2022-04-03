package storage

import (
	"testing"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/matryer/is"
)

const dbfile = "./../../mygopher.db"

func TestSaveNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	tagsStore := Tag{conn}
	tagsSvc := tagging.New(tagsStore)

	tg1, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg2, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	want := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg1, tg2})
	err = newsStore.Save(*want)
	if err != nil {
		t.Fatal(err.Error())
	}

	is := is.New(t)

	got, err := newsStore.GetById(want.Post.ID.String())
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Post.Status, want.Post.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
	is.Equal(len(got.Tags), len(want.Tags))
}

func TestUpdateNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	tagsStore := Tag{conn}
	tagsSvc := tagging.New(tagsStore)

	tg1, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg2, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg3, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	news := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg1, tg2})
	err = newsStore.Save(*news)
	if err != nil {
		t.Fatal(err.Error())
	}

	wantTitle := gofakeit.Sentence(10)
	wantTag := []domain.Tags{tg3}
	wantStatus := domain.Publish

	news.ChangeTitle(wantTitle)
	news.ChangeTags(wantTag)
	news.ChangeStatus(wantStatus)

	err = newsStore.Update(*news)
	if err != nil {
		t.Fatal(err.Error())
	}

	is := is.New(t)

	got, err := newsStore.GetById(news.Post.ID.String())
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, wantTitle)
	is.Equal(got.Post.Status, wantStatus)
	is.Equal(len(got.Tags), len(wantTag))
}

func TestGetByID(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	tagsStore := Tag{conn}
	tagsSvc := tagging.New(tagsStore)

	tg1, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg2, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	want := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg1, tg2})
	err = newsStore.Save(*want)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetById(want.Post.ID.String())

	is := is.New(t)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Post.Status, want.Post.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
	is.Equal(len(got.Tags), len(want.Tags))
}

func TestGetAllNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	tagsStore := Tag{conn}
	tagsSvc := tagging.New(tagsStore)

	tg1, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg2, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	wantNews1 := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg1})
	err = newsStore.Save(*wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	wantNews2 := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg2})
	err = newsStore.Save(*wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAll()
	is := is.New(t)
	is.NoErr(err)
	is.Equal(len(got), 2)

	is.Equal(got[0].Post.Title, wantNews1.Post.Title)
	is.Equal(got[0].Post.Status, wantNews1.Post.Status)
	is.Equal(got[0].Post.Body, wantNews1.Post.Body)
	is.Equal(got[0].Slug, wantNews1.Slug)
	is.Equal(len(got[0].Tags), len(wantNews1.Tags))

	is.Equal(got[1].Post.Title, wantNews2.Post.Title)
	is.Equal(got[1].Post.Status, wantNews2.Post.Status)
	is.Equal(got[1].Post.Body, wantNews2.Post.Body)
	is.Equal(got[1].Slug, wantNews2.Slug)
	is.Equal(len(got[1].Tags), len(wantNews2.Tags))
}

func TestDeleteNews(t *testing.T) {
	conn := Run(dbfile,true)
	newsStore := News{conn}
	tagsStore := Tag{conn}
	tagsSvc := tagging.New(tagsStore)

	tg1, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	tg2, err := tagsSvc.Create(gofakeit.Noun())
	if err != nil {
		t.Fatal(err.Error())
	}

	news := news.New(gofakeit.Sentence(10), gofakeit.Sentence(100), domain.Draft, []domain.Tags{tg1, tg2})
	err = newsStore.Save(*news)
	if err != nil {
		t.Fatal(err.Error())
	}

	is := is.New(t)

	err = newsStore.Delete(news.Post.ID)
	is.NoErr(err)

	_, err = newsStore.GetById(news.Post.ID.String())
	is.True(err != nil)
}
