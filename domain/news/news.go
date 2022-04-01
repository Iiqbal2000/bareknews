package news

import (
	"errors"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
	"github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidTitle = errors.New("invalid title")
)

// Agregates
type News struct {
	Post *bareknews.Posts
	Slug bareknews.Slug
	Tags []bareknews.Tags
}

func New(title, body string, status bareknews.Status, tags []bareknews.Tags) *News {
	post := &bareknews.Posts{
		ID: uuid.New(),
		Status: status,
		Title: title,
		Body: body,
	}

	return &News{
		Post: post,
		Slug: bareknews.NewSlug(post.Title),
		Tags: tags,
	}
}

func (n News) Validate() error {
	return validation.ValidateStruct(&n, 
		validation.Field(&n.Post.Title, validation.Required, validation.Length(5, 50)),
		validation.Field(&n.Post.Body, validation.Required),
		validation.Field(&n.Post.Status,
			validation.Required,
			validation.In(bareknews.Publish, bareknews.Draft, bareknews.Deleted),
		),
	)
}

func (n *News) ChangeTitle(newTitle string) {
	n.Post.Title = newTitle
	n.changeSlug()
}

func (n *News) changeSlug() {
	n.Slug = bareknews.NewSlug(n.Post.Title)
}

func (n *News) ChangeStatus(newStatus bareknews.Status) {
	n.Post.Status = newStatus
}

func (n *News) ChangeBody(newBody string) {
	n.Post.Body = newBody
}

func (n *News) ChangeTags(newTags []bareknews.Tags) {
	n.Tags = newTags
}