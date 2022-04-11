package news

import (
	"errors"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidTitle = errors.New("invalid title")
)

// News is an aggregate that represents a piece of news.
type News struct {
	Post   domain.Post
	Status domain.Status `json:"status"`
	Slug   domain.Slug   `json:"slug"`
	TagsID []uuid.UUID
}

func New(title, body string, status domain.Status, tags []uuid.UUID) *News {
	post := domain.Post{
		ID:    uuid.New(),
		Title: title,
		Body:  body,
	}

	return &News{
		Post:   post,
		Status: status,
		Slug:   domain.NewSlug(post.Title),
		TagsID: tags,
	}
}

func (n News) Validate() error {
	if err := n.Post.Validate(); err != nil {
		return err
	}

	if err := n.Status.Validate(); err != nil {
		return err
	}

	return nil
}

func (n *News) ChangeTitle(newTitle string) {
	n.Post.Title = newTitle
	n.changeSlug()
}

func (n *News) changeSlug() {
	n.Slug = domain.NewSlug(n.Post.Title)
}

func (n *News) ChangeStatus(newStatus domain.Status) {
	n.Status = newStatus
}

func (n *News) ChangeBody(newBody string) {
	n.Post.Body = newBody
}

func (n *News) ChangeTags(newTags []uuid.UUID) {
	n.TagsID = newTags
}
