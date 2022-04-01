package news_test

// import (
// 	"testing"

// 	"github.com/Iiqbal2000/bareknews"
// 	"github.com/Iiqbal2000/bareknews/domain/news"
// 	"github.com/matryer/is"
// )

// func TestCreateNews(t *testing.T) {
// 	title := "Test 1"
// 	slug := bareknews.Slug("test-1")
// 	body := "testing"
// 	tag1, _ := bareknews.NewTags("tag1")
// 	tag2, _ := bareknews.NewTags("tag2")
// 	news, err := news.New(title, body, []bareknews.Tags{*tag1, *tag2})
	
// 	is := is.New(t)
// 	is.True(err == nil)
// 	is.Equal(news.Post.Title, title)
// 	is.Equal(news.Post.Body, body)
// 	is.Equal(news.Post.Status, bareknews.Draft)
// 	is.Equal(news.Slug, slug)
// }

// func TestChangeTitle(t *testing.T) {
// 	tag1, _ := bareknews.NewTags("tag1")
// 	tag2, _ := bareknews.NewTags("tag2")
// 	news, err := news.New("Test 1", "testing", []bareknews.Tags{*tag1, *tag2})
	
// 	is := is.New(t)
// 	is.True(err == nil)

// 	news.ChangeTitle("Test 2")
// 	is.Equal(news.Post.Title, "Test 2")
// 	is.Equal(news.Slug, bareknews.Slug("test-2"))
// }

// func TestChangeBody(t *testing.T) {
// 	tag1, _ := bareknews.NewTags("tag1")
// 	tag2, _ := bareknews.NewTags("tag2")
// 	news, err := news.New("Test 1", "testing", []bareknews.Tags{*tag1, *tag2})
	
// 	is := is.New(t)
// 	is.True(err == nil)

// 	news.ChangeBody("Changing the body")
// 	is.Equal(news.Post.Body, "Changing the body")
// }

// func TestChangeStatus(t *testing.T) {
// 	tag1, _ := bareknews.NewTags("tag1")
// 	tag2, _ := bareknews.NewTags("tag2")
// 	news, err := news.New("Test 1", "testing", []bareknews.Tags{*tag1, *tag2})
	
// 	is := is.New(t)
// 	is.True(err == nil)

// 	news.ChangeStatus(int(bareknews.Publish))
// 	is.Equal(news.Post.Status, bareknews.Publish)
// }