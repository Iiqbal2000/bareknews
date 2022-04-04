package news

import (
	"errors"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/google/uuid"
	"github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidTitle = errors.New("invalid title")
)

// Agregates
type News struct {
	Post domain.Posts
	Slug domain.Slug	`json:"slug"`
	Tags []domain.Tags
}

func New(title, body string, status domain.Status, tags []domain.Tags) *News {
	post := domain.Posts{
		ID: uuid.New(),
		Status: status,
		Title: title,
		Body: body,
	}

	return &News{
		Post: post,
		Slug: domain.NewSlug(post.Title),
		Tags: tags,
	}
}

func (n News) Validate() error {
	if err := n.Post.Status.Validate(); err != nil {
		return err
	}
	
	return validation.ValidateStruct(&n, 
		validation.Field(&n.Post.Title, validation.Required, validation.Length(5, 50)),
		validation.Field(&n.Post.Body, validation.Required),
	)
}

func (n *News) ChangeTitle(newTitle string) {
	n.Post.Title = newTitle
	n.changeSlug()
}

func (n *News) changeSlug() {
	n.Slug = domain.NewSlug(n.Post.Title)
}

func (n *News) ChangeStatus(newStatus domain.Status) {
	n.Post.Status = newStatus
}

func (n *News) ChangeBody(newBody string) {
	n.Post.Body = newBody
}

func (n *News) ChangeTags(newTags []domain.Tags) {
	n.Tags = newTags
}