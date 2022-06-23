package db

import (
	"context"
	"database/sql"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	conn *sql.DB
}

func CreateStore(conn *sql.DB) Store {
	return Store{conn: conn}
}

func (t Store) Save(ctx context.Context, tag tags.Tags) error {
	query, args := sqlbuilder.InsertInto("tags").
		Cols("id", "name", "slug").
		Values(tag.Label.ID, tag.Label.Name, tag.Slug).
		Build()

	_, err := t.conn.ExecContext(ctx, query, args...)
	if err != nil {
		if possibleErr, ok := err.(sqlite3.Error); ok {
			if possibleErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return bareknews.ErrDataAlreadyExist
			}
		}

		return errors.Wrap(err, "storage.tags.save")
	}

	return nil
}

func (t Store) Update(ctx context.Context, tag tags.Tags) error {
	builder := sqlbuilder.NewUpdateBuilder()
	builder.Update("tags")
	builder.Set(
		builder.Assign("name", tag.Label.Name),
		builder.Assign("slug", tag.Slug),
	)
	builder.Where(builder.Equal("id", tag.Label.ID.String()))
	query, args := builder.Build()
	_, err := t.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "storage.tags.update")
	}

	return nil
}

func (t Store) Delete(ctx context.Context, id uuid.UUID) error {
	d := sqlbuilder.NewDeleteBuilder()
	d.DeleteFrom("tags")
	d.Where(d.Equal("id", id))

	query, args := d.Build()

	_, err := t.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "storage.tags.delete")
	}

	return nil
}

func (t Store) GetById(ctx context.Context, id uuid.UUID) (*tags.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.Equal("id", id))
	query, args := builder.Build()
	row := t.conn.QueryRowContext(ctx, query, args...)

	label := bareknews.Label{}
	var slug bareknews.Slug
	err := row.Scan(&label.ID, &label.Name, &slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &tags.Tags{}, sql.ErrNoRows
		} else {
			return &tags.Tags{}, errors.Wrap(err, "storage.tags.getById")
		}
	}

	tag := &tags.Tags{
		Label: label,
		Slug:  slug,
	}

	return tag, nil
}

func (t Store) GetByIds(ctx context.Context, ids []uuid.UUID) ([]tags.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	idstr := make([]string, 0)

	for _, elem := range ids {
		idstr = append(idstr, elem.String())
	}

	listMark := sqlbuilder.List(idstr)

	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.In("id", listMark))
	query, args := builder.Build()

	rows, err := t.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []tags.Tags{}, errors.Wrap(err, "storage.tags.getByIds")
	}

	defer rows.Close()

	results := make([]tags.Tags, 0)

	for rows.Next() {
		label := bareknews.Label{}
		var slug bareknews.Slug
		err := rows.Scan(&label.ID, &label.Name, &slug)
		if err != nil {
			return []tags.Tags{}, errors.Wrap(err, "storage.tags.getByIds")
		}

		results = append(results, tags.Tags{
			Label: label,
			Slug:  slug,
		})
	}

	if rows.Err() != nil {
		return []tags.Tags{}, errors.Wrap(err, "storage.tags.getByIds")
	}

	return results, nil
}

func (t Store) GetAll(ctx context.Context) ([]tags.Tags, error) {
	query, args := sqlbuilder.Select("id", "name", "slug").
		From("tags").
		Build()

	rows, err := t.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []tags.Tags{}, errors.Wrap(err, "storage.tags.getAll")
	}

	defer rows.Close()

	results := make([]tags.Tags, 0)

	for rows.Next() {
		label := bareknews.Label{}
		var slug bareknews.Slug
		err := rows.Scan(&label.ID, &label.Name, &slug)
		if err != nil {
			return []tags.Tags{}, errors.Wrap(err, "storage.tags.getAll")
		}

		results = append(results, tags.Tags{
			Label: label,
			Slug:  slug,
		})
	}

	if rows.Err() != nil {
		return []tags.Tags{}, errors.Wrap(err, "storage.tags.getAll")
	}

	return results, nil
}

func (t Store) Count(ctx context.Context, id uuid.UUID) (int, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select(builder.As("COUNT(id)", "c"))
	builder.From("tags")
	builder.Where(builder.Equal("id", id))
	query, args := builder.Build()
	row := t.conn.QueryRowContext(ctx, query, args...)

	var c int
	err := row.Scan(&c)
	if err != nil {
		return c, errors.Wrap(err, "storage.tags.count")
	}

	if c == 0 {
		return c, sql.ErrNoRows
	}

	return c, nil
}

func (t Store) GetByNames(ctx context.Context, names ...string) ([]tags.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	listMark := sqlbuilder.List(names)
	builder.Select("id", "name", "slug")
	builder.From("tags")
	builder.Where(builder.In("name", listMark))
	query, args := builder.Build()

	rows, err := t.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []tags.Tags{}, errors.Wrap(err, "storage.tags.getByNames")
	}

	defer rows.Close()

	results := make([]tags.Tags, 0)

	for rows.Next() {
		label := bareknews.Label{}
		var slug bareknews.Slug
		err := rows.Scan(&label.ID, &label.Name, &slug)
		if err != nil {
			return []tags.Tags{}, errors.Wrap(err, "storage.tags.getByNames")
		}

		results = append(results, tags.Tags{
			Label: label,
			Slug:  slug,
		})
	}

	return results, nil
}

func (t Store) GetByName(ctx context.Context, name string) (tags.Tags, error) {
	queryBuilder := sqlbuilder.NewSelectBuilder()
	queryBuilder.Select("id", "name", "slug")
	queryBuilder.From("tags")
	queryBuilder.Where(queryBuilder.Equal("name", name))
	
	query, args := queryBuilder.Build()
	row := t.conn.QueryRowContext(ctx, query, args...)

	label := bareknews.Label{}
	var slug bareknews.Slug
	
	err := row.Scan(&label.ID, &label.Name, &slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tags.Tags{}, sql.ErrNoRows
		} else {
			return tags.Tags{}, errors.Wrap(err, "storage.tags.GetByName")
		}
	}

	tag := tags.Tags{
		Label: label,
		Slug:  slug,
	}

	return tag, nil
}

