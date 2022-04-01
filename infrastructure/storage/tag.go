package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Iiqbal2000/bareknews"
	"github.com/huandu/go-sqlbuilder"
)

type TagStorage struct {
	Conn *sql.DB
}

var tags = sqlbuilder.NewStruct(new(bareknews.Tags))

func (t TagStorage) Save(tag bareknews.Tags) error {
	query, args := sqlbuilder.InsertInto("tags").
		Cols("id", "name", "slug").
		Values(tag.ID, tag.Name, tag.Slug).
		Build()

	_, err := t.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failure when inserting a tag: %s", err.Error())
	}

	return nil
}

func (t TagStorage) Update(tag bareknews.Tags) error {
	builder := sqlbuilder.NewUpdateBuilder()
	builder.Update("tags")
	builder.Set(
		builder.Assign("name", tag.Name),
		builder.Assign("slug", tag.Slug),
	)
	builder.Where(builder.Equal("id", tag.ID))
	query, args := builder.Build()
	_, err := t.Conn.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not found")	
		} else {
			return fmt.Errorf("failure when updating a tag: %s", err.Error())
		}
	}

	return nil
}

func (t TagStorage) GetById(id string) (*bareknews.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.Equal("id", id))
	query, args := builder.Build()
	row := t.Conn.QueryRow(query, args...)
	
	tag := &bareknews.Tags{}
	err := row.Scan(tags.Addr(tag)...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &bareknews.Tags{}, fmt.Errorf("not found")	
		} else {
			return &bareknews.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}
	}

	return tag, nil
}

func (t TagStorage) GetByNewsId(id string) ([]bareknews.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Distinct()
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Join("news_tags", builder.Equal("news_tags.newsID", id))
	query, args := builder.Build()
	rows, err := t.Conn.Query(query, args...)
	if err != nil {
		return []bareknews.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
	}

	defer rows.Close()

	results := make([]bareknews.Tags, 0)

	for rows.Next() {
		t := bareknews.Tags{}
		err = rows.Scan(tags.Addr(&t)...)
		if err != nil {
			return []bareknews.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}
		results = append(results, t)
	}

	return results, nil
}

func (t TagStorage) GetByName(name string) (bareknews.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.Equal("name", name))
	query, args := builder.Build()
	row := t.Conn.QueryRow(query, args...)
	
	tag := &bareknews.Tags{}
	err := row.Scan(tags.Addr(tag)...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bareknews.Tags{}, fmt.Errorf("not found")	
		} else {
			return bareknews.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}
	}

	return *tag, nil
}

func (t TagStorage) GetAll() ([]bareknews.Tags, error) {
	query, args := sqlbuilder.Select("id", "name", "slug").
		From("tags").
		Build()

	rows, err := t.Conn.Query(query, args...)
	if err != nil {
		return []bareknews.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
	}

	defer rows.Close()

	results := make([]bareknews.Tags, 0)

	for rows.Next() {
		t := bareknews.Tags{}
		err = rows.Scan(tags.Addr(&t)...)
		if err != nil {
			return []bareknews.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}
		results = append(results, t)
	}

	return results, nil
}

func (t TagStorage) GetByNames(names ...string) ([]bareknews.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	listMark := sqlbuilder.List(names)
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.In("name", listMark))
	query, args := builder.Build()

	rows, err := t.Conn.Query(query, args...)
	if err != nil {
		return []bareknews.Tags{}, fmt.Errorf("failure when querying tags: %s", err.Error())
	}

	defer rows.Close()

	results := make([]bareknews.Tags, 0)

	for rows.Next() {
		t := bareknews.Tags{}
		err = rows.Scan(tags.Addr(&t)...)
		if err != nil {
			return []bareknews.Tags{}, fmt.Errorf("failure when querying tags in the iteration: %s", err.Error())
		}
		results = append(results, t)
	}

	return results, nil
}

func (t TagStorage) Delete(id string) error {
	d := sqlbuilder.NewDeleteBuilder()
	d.DeleteFrom("tags")
	d.Where(d.Equal("id", id))
	sql, args := d.Build()
	_, err := t.Conn.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failure when deleting a tag: %s", err.Error())
	}

	return nil
}
