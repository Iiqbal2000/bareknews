package storage

import (
	"context"
	"testing"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestSaveNews(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	want := news.New(title, body, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *want)
	is.NoErr(err)

	got, err := newsStore.GetById(context.TODO(), want.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Status, want.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
}

func TestUpdateNews(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}
	is := is.New(t)
	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	news := news.New(title, body, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *news)
	is.NoErr(err)

	wantTitle := "news 2"
	wantStatus := domain.Publish

	news.ChangeTitle(wantTitle)
	news.ChangeStatus(wantStatus)

	err = newsStore.Update(context.TODO(), *news)
	is.NoErr(err)

	got, err := newsStore.GetById(context.TODO(), news.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, wantTitle)
	is.Equal(got.Status, wantStatus)
}

func TestDeleteNews(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	news := news.New(title, body, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *news)
	is.NoErr(err)

	err = newsStore.Delete(context.TODO(), news.Post.ID)
	is.NoErr(err)

	_, err = newsStore.GetById(context.TODO(), news.Post.ID)
	is.True(err != nil)
}

func TestGetByID(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}
	is := is.New(t)

	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	want := news.New(title, body, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *want)
	is.NoErr(err)

	got, err := newsStore.GetById(context.TODO(), want.Post.ID)

	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Status, want.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
	is.Equal(len(got.TagsID), 1)
}

func TestGetAllNews(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.New(title1, body1, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.New(title2, body2, domain.Draft, []uuid.UUID{tgId})
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAll(context.TODO())
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

func TestGetAllByTopic(t *testing.T) {
	conn := Run(":memory:", true)
	newsStore := News{conn}

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.New(title1, body1, domain.Draft, []uuid.UUID{tgId})
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.New(title2, body2, domain.Draft, []uuid.UUID{tgId})
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAllByTopic(context.TODO(), tgId)
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