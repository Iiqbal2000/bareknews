package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/Iiqbal2000/bareknews/news/db"
	"github.com/Iiqbal2000/bareknews/pkg/sqlite3"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestSaveNews(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)
	is := is.New(t)

	tgIds := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	want := news.Create(title, body, bareknews.Draft, tgIds, time.Now().Unix())
	err := newsStore.Save(context.TODO(), *want)
	is.NoErr(err)

	got, err := newsStore.GetById(context.TODO(), want.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, want.Post.Title)
	is.Equal(got.Status, want.Status)
	is.Equal(got.Post.Body, want.Post.Body)
	is.Equal(got.Slug, want.Slug)
	is.Equal(len(got.TagsID), len(tgIds))
	is.Equal(got.DateCreated, want.DateCreated)
	is.Equal(got.DateUpdated, want.DateUpdated)
}

func TestUpdateNews(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)
	is := is.New(t)
	tgIds := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."
	unixTime := time.Now().Unix()

	news := news.Create(title, body, bareknews.Draft, tgIds, unixTime)
	err := newsStore.Save(context.TODO(), *news)
	is.NoErr(err)

	wantTitle := "news 2"
	wantStatus := bareknews.Publish
	wantTags := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()}

	news.ChangeTitle(wantTitle)
	news.ChangeStatus(wantStatus)
	news.ChangeTags(wantTags)
	news.ChangeDateUpdated(time.Now().Unix())

	err = newsStore.Update(context.TODO(), *news)
	is.NoErr(err)

	got, err := newsStore.GetById(context.TODO(), news.Post.ID)
	is.NoErr(err)
	is.True(got != nil)
	is.Equal(got.Post.Title, wantTitle)
	is.Equal(got.Status, wantStatus)
	is.Equal(len(got.TagsID), len(wantTags))
}

func TestDeleteNews(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)
	is := is.New(t)

	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	news := news.Create(title, body, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err := newsStore.Save(context.TODO(), *news)
	is.NoErr(err)

	err = newsStore.Delete(context.TODO(), news.Post.ID)
	is.NoErr(err)

	_, err = newsStore.GetById(context.TODO(), news.Post.ID)
	is.True(err != nil)
}

func TestGetByID(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)
	is := is.New(t)

	tgId := uuid.New()

	title := "news 1"
	body := "Struct fields can also use tags to more specifically generate data for that field type."

	want := news.Create(title, body, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
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
	is.True(got.DateCreated != 0)
	is.True(got.DateUpdated != 0)
}

func TestGetAllNews(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: "./../../bareknews.db", DropTableFirst: true})
	newsStore := db.CreateStore(conn)

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.Create(title1, body1, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.Create(title2, body2, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAll(context.TODO(), 0, 2)
	is := is.New(t)
	is.NoErr(err)
	is.Equal(len(got), 2)

	is.Equal(got[0].Post.Title, wantNews1.Post.Title)
	is.Equal(got[0].Status, wantNews1.Status)
	is.Equal(got[0].Post.Body, wantNews1.Post.Body)
	is.Equal(got[0].Slug, wantNews1.Slug)
	is.Equal(len(got[0].TagsID), 1)
	is.True(got[0].DateCreated != 0)
	is.True(got[0].DateUpdated != 0)

	is.Equal(got[1].Post.Title, wantNews2.Post.Title)
	is.Equal(got[1].Status, wantNews2.Status)
	is.Equal(got[1].Post.Body, wantNews2.Post.Body)
	is.Equal(got[1].Slug, wantNews2.Slug)
	is.Equal(len(got[1].TagsID), 1)
	is.True(got[1].DateCreated != 0)
	is.True(got[1].DateUpdated != 0)
}

func TestGetAllByTopic(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: "./../../bareknews.db", DropTableFirst: true})
	newsStore := db.CreateStore(conn)

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.Create(title1, body1, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.Create(title2, body2, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAllByTopic(context.TODO(), tgId, 0, 2)
	is := is.New(t)
	is.NoErr(err)
	is.Equal(len(got), 2)
	is.True(got[0].DateCreated != 0)
	is.True(got[0].DateUpdated != 0)
	is.True(got[1].DateCreated != 0)
	is.True(got[1].DateUpdated != 0)
}

func TestGetAllByStatus(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.Create(title1, body1, bareknews.Publish, []uuid.UUID{tgId}, time.Now().Unix())
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.Create(title2, body2, bareknews.Draft, []uuid.UUID{tgId}, time.Now().Unix())
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	got, err := newsStore.GetAllByStatus(context.TODO(), bareknews.Publish, 0, 2)
	is := is.New(t)
	is.NoErr(err)
	is.Equal(len(got), 1)

	is.Equal(got[0].Post.Title, wantNews1.Post.Title)
	is.Equal(got[0].Status, wantNews1.Status)
	is.Equal(got[0].Post.Body, wantNews1.Post.Body)
	is.Equal(got[0].Slug, wantNews1.Slug)
	is.True(got[0].DateCreated != 0)
	is.True(got[0].DateUpdated != 0)
	is.Equal(len(got[0].TagsID), 1)
}

func TestGetAllByPaginationNoCursor(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.Create(title1, body1, bareknews.Publish, []uuid.UUID{tgId}, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.Create(title2, body2, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	title3 := "news 3"
	body3 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews3 := news.Create(title3, body3, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews3)
	if err != nil {
		t.Fatal(err.Error())
	}

	title4 := "news 4"
	body4 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews4 := news.Create(title4, body4, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews4)
	if err != nil {
		t.Fatal(err.Error())
	}

	is := is.New(t)

	got, err := newsStore.GetAllByPagination(context.TODO(), 0, 2)
	is.NoErr(err)

	is.Equal(len(got), 2)
}

func TestGetAllByPaginationWithCursor(t *testing.T) {
	conn, _ := sqlite3.Run(sqlite3.Config{URI: ":memory:", DropTableFirst: true})
	newsStore := db.CreateStore(conn)

	tgId := uuid.New()

	title1 := "news 1"
	body1 := "Struct fields can also use tags to more specifically generate data for that field type."

	wantNews1 := news.Create(title1, body1, bareknews.Publish, []uuid.UUID{tgId}, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err := newsStore.Save(context.TODO(), *wantNews1)
	if err != nil {
		t.Fatal(err.Error())
	}

	title2 := "news 2"
	body2 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews2 := news.Create(title2, body2, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews2)
	if err != nil {
		t.Fatal(err.Error())
	}

	title3 := "news 3"
	body3 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews3 := news.Create(title3, body3, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews3)
	if err != nil {
		t.Fatal(err.Error())
	}

	title4 := "news 4"
	body4 := "Lorem Ipsum is simply dummy text of the printing and typesetting industry."

	wantNews4 := news.Create(title4, body4, bareknews.Draft, []uuid.UUID{tgId}, time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).Unix())
	err = newsStore.Save(context.TODO(), *wantNews4)
	if err != nil {
		t.Fatal(err.Error())
	}

	is := is.New(t)

	got, err := newsStore.GetAllByPagination(context.TODO(), wantNews3.DateCreated, 2)
	is.NoErr(err)

	is.Equal(len(got), 2)
}