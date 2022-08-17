package news_test

import (
	"testing"
	"time"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreateNews(t *testing.T) {
	t.Run("success creates a news", func(t *testing.T) {
		title := "Test 1"
		slug := bareknews.Slug("test-1")
		body := "testing"
		tagid := uuid.New()
		news := news.Create(title, body, bareknews.Draft, []uuid.UUID{tagid}, time.Now().Unix())

		err := news.Validate()

		is := is.New(t)
		is.NoErr(err)
		is.Equal(news.Post.Title, title)
		is.Equal(news.Post.Body, body)
		is.Equal(news.Status, bareknews.Draft)
		is.Equal(news.Slug, slug)
	})

	t.Run("with invalid title should return an error", func(t *testing.T) {
		body := "testing"
		tagid := uuid.New()
		news := news.Create("", body, bareknews.Draft, []uuid.UUID{tagid}, time.Now().Unix())

		err := news.Validate()

		is := is.New(t)
		is.True(err != nil)
	})

	t.Run("with invalid status should return an error", func(t *testing.T) {
		body := "testing"
		tagid := uuid.New()
		news := news.Create("invalid status", body, "", []uuid.UUID{tagid}, time.Now().Unix())

		err := news.Validate()

		is := is.New(t)
		is.True(err != nil)
	})
}

func TestChangeTitle(t *testing.T) {
	tagId := uuid.New()
	news := news.Create("Test 1", "testing", bareknews.Draft, []uuid.UUID{tagId}, time.Now().Unix())
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)

	news.ChangeTitle("Test 2")
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.Post.Title, "Test 2")
	is.Equal(news.Slug, bareknews.Slug("test-2"))
}

func TestChangeBody(t *testing.T) {
	tagId := uuid.New()
	news := news.Create("Test 1", "testing", bareknews.Draft, []uuid.UUID{tagId}, time.Now().Unix())
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)

	news.ChangeBody("Changing the body")
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.Post.Body, "Changing the body")
}

func TestChangeStatus(t *testing.T) {
	tagId := uuid.New()
	news := news.Create("Test 1", "testing", bareknews.Draft, []uuid.UUID{tagId}, time.Now().Unix())
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)

	news.ChangeStatus(bareknews.Publish)
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.Status, bareknews.Publish)
}

func TestChangeTags(t *testing.T) {
	tagId := uuid.New()
	news := news.Create("Test 1", "testing", bareknews.Draft, []uuid.UUID{tagId}, time.Now().Unix())
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)

	newTag := uuid.New()

	news.ChangeTags([]uuid.UUID{newTag})
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.TagsID[0], newTag)
}