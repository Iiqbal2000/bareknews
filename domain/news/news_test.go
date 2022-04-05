package news_test

import (
	"testing"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestCreateNews(t *testing.T) {
	t.Run("success creates a news", func(t *testing.T) {
		title := "Test 1"
		slug := domain.Slug("test-1")
		body := "testing"
		tagid := uuid.New()
		news := news.New(title, body, domain.Draft, []uuid.UUID{tagid})

		err := news.Validate()

		is := is.New(t)
		is.NoErr(err)
		is.Equal(news.Post.Title, title)
		is.Equal(news.Post.Body, body)
		is.Equal(news.Status, domain.Draft)
		is.Equal(news.Slug, slug)
	})

	t.Run("with invalid title should return an error", func(t *testing.T) {
		body := "testing"
		tagid := uuid.New()
		news := news.New("", body, domain.Draft, []uuid.UUID{tagid})

		err := news.Validate()

		is := is.New(t)
		is.True(err != nil)
	})

	t.Run("with invalid status should return an error", func(t *testing.T) {
		body := "testing"
		tagid := uuid.New()
		news := news.New("invalid status", body, "", []uuid.UUID{tagid})

		err := news.Validate()

		is := is.New(t)
		is.True(err != nil)
	})
}

func TestChangeTitle(t *testing.T) {
	tagId := uuid.New()
	news := news.New("Test 1", "testing", domain.Draft, []uuid.UUID{tagId})
	err := news.Validate()
	
	is := is.New(t)
	is.NoErr(err)

	news.ChangeTitle("Test 2")
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.Post.Title, "Test 2")
	is.Equal(news.Slug, domain.Slug("test-2"))
}

func TestChangeBody(t *testing.T) {
	tagId := uuid.New()
	news := news.New("Test 1", "testing", domain.Draft, []uuid.UUID{tagId})
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
	news := news.New("Test 1", "testing", domain.Draft, []uuid.UUID{tagId})
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)

	news.ChangeStatus(domain.Publish)
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.Status, domain.Publish)
}

func TestChangeTags(t *testing.T) {
	tagId := uuid.New()
	news := news.New("Test 1", "testing", domain.Draft, []uuid.UUID{tagId})
	err := news.Validate()

	is := is.New(t)
	is.NoErr(err)
	
	newTag := uuid.New()

	news.ChangeTags([]uuid.UUID{newTag})
	err = news.Validate()
	is.NoErr(err)
	is.Equal(news.TagsID[0], newTag)
}