package news

import (
	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

// News is an aggregate that represents a piece of news.
type News struct {
	Post   bareknews.Post
	Status bareknews.Status
	Slug   bareknews.Slug
	TagsID []uuid.UUID
	DateCreated int64
	DateUpdated int64
}

func Create(title, body string, status bareknews.Status, tags []uuid.UUID, timeNowUnix int64) *News {
	post := bareknews.Post{
		ID:    uuid.New(),
		Title: title,
		Body:  body,
	}

	return &News{
		Post:   post,
		Status: status,
		Slug:   bareknews.NewSlug(post.Title),
		TagsID: tags,
		DateCreated: timeNowUnix,
		DateUpdated: timeNowUnix,
	}
}

func Update(id uuid.UUID, title, body string, status bareknews.Status, tags []uuid.UUID, dateCreated, dateUpdated int64) *News {
	post := bareknews.Post{
		ID:    id,
		Title: title,
		Body:  body,
	}

	return &News{
		Post:   post,
		Status: status,
		Slug:   bareknews.NewSlug(post.Title),
		TagsID: tags,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
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
	n.Slug = bareknews.NewSlug(n.Post.Title)
}

func (n *News) ChangeStatus(newStatus bareknews.Status) {
	n.Status = newStatus
}

func (n *News) ChangeBody(newBody string) {
	n.Post.Body = newBody
}

func (n *News) ChangeTags(newTags []uuid.UUID) {
	n.TagsID = newTags
}

func (n *News) ChangeDateUpdated(timeNowUnix int64) {
	n.DateUpdated = timeNowUnix
}