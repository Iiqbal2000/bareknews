package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/tags"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type Tag struct {
	Conn *sql.DB
}

func (t Tag) Save(tag tags.Tags) error {
	query, args := sqlbuilder.InsertInto("tags").
		Cols("id", "name", "slug").
		Values(tag.Label.ID, tag.Label.Name, tag.Slug).
		Build()

	_, err := t.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failure when inserting a tag: %s", err.Error())
	}

	return nil
}

func (t Tag) Update(tag tags.Tags) error {
	builder := sqlbuilder.NewUpdateBuilder()
	builder.Update("tags")
	builder.Set(
		builder.Assign("name", tag.Label.Name),
		builder.Assign("slug", tag.Slug),
	)
	builder.Where(builder.Equal("id", tag.Label.ID.String()))
	query, args := builder.Build()
	_, err := t.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failure when updating a tag: %s", err.Error())
	}

	return nil
}

func (t Tag) Delete(id uuid.UUID) error {
	d := sqlbuilder.NewDeleteBuilder()
	d.DeleteFrom("tags")
	d.Where(d.Equal("id", id))
	query, args := d.Build()
	_, err := t.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failure when deleting a tag: %s", err.Error())
	}

	return nil
}

func (t Tag) GetById(id uuid.UUID) (*tags.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.Equal("id", id))
	query, args := builder.Build()
	row := t.Conn.QueryRow(query, args...)

	label := domain.Label{}
	var slug domain.Slug
	err := row.Scan(&label.ID, &label.Name, &slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &tags.Tags{}, sql.ErrNoRows
		} else {
			return &tags.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
		}
	}

	tag := &tags.Tags{
		Label: label,
		Slug:  slug,
	}

	return tag, nil
}

func (t Tag) GetAll() ([]tags.Tags, error) {
	query, args := sqlbuilder.Select("id", "name", "slug").
		From("tags").
		Build()

	rows, err := t.Conn.Query(query, args...)
	if err != nil {
		return []tags.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
	}

	defer rows.Close()

	results := make([]tags.Tags, 0)

	for rows.Next() {
		label := domain.Label{}
		var slug domain.Slug
		err := rows.Scan(&label.ID, &label.Name, &slug)
		if err != nil {
			return []tags.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}

		results = append(results, tags.Tags{
			Label: label,
			Slug:  slug,
		})
	}

	if rows.Err() != nil {
		return []tags.Tags{}, fmt.Errorf("tags are incomplete")
	}

	return results, nil
}

// func (t Tag) GetByNewsId(id string) ([]tags.Tags, error) {
// 	builder := sqlbuilder.NewSelectBuilder()
// 	builder.Distinct()
// 	builder.Select("id", "name", "slug")
// 	builder.From("tags")
// 	builder.Join("news_tags", builder.Equal("news_tags.newsID", id))
// 	query, args := builder.Build()
// 	rows, err := t.Conn.Query(query, args...)
// 	if err != nil {
// 		return []tags.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
// 	}

// 	defer rows.Close()

// 	results := make([]tags.Tags, 0)

// 	for rows.Next() {
// 		t := tags.Tags{}
// 		err = rows.Scan(tags.Addr(&t)...)
// 		if err != nil {
// 			return []tags.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
// 		}
// 		results = append(results, t)
// 	}

// 	return results, nil
// }

// func (t Tag) GetByName(name string) (tags.Tags, error) {
// 	builder := sqlbuilder.NewSelectBuilder()
// 	builder.Select("id", "name", "slug")
// 	builder.From("tags")
// 	builder.Where(builder.Equal("name", name))
// 	query, args := builder.Build()
// 	row := t.Conn.QueryRow(query, args...)

// 	tag := &tags.Tags{}
// 	err := row.Scan(tags.Addr(tag)...)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return tags.Tags{}, fmt.Errorf("not found")
// 		} else {
// 			return tags.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
// 		}
// 	}

// 	return *tag, nil
// }

// func (t Tag) GetByNames(names ...string) ([]tags.Tags, error) {
// 	builder := sqlbuilder.NewSelectBuilder()
// 	listMark := sqlbuilder.List(names)
// 	builder.Select("id", "name", "slug")
// 	builder.From("tags")
// 	builder.Where(builder.In("name", listMark))
// 	query, args := builder.Build()

// 	rows, err := t.Conn.Query(query, args...)
// 	if err != nil {
// 		return []tags.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
// 	}

// 	defer rows.Close()

// 	results := make([]tags.Tags, 0)

// 	for rows.Next() {
// 		t := tags.Tags{}
// 		err = rows.Scan(tags.Addr(&t)...)
// 		if err != nil {
// 			return []tags.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
// 		}
// 		results = append(results, t)
// 	}

// 	return results, nil
// }
