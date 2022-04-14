package storage

import (
	"database/sql"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type News struct {
	Conn *sql.DB
}

func (s News) Save(n news.News) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "storage.news.save")
	}

	defer tx.Rollback()

	builder := sqlbuilder.NewInsertBuilder()
	builder.InsertInto("news")
	builder.Cols("id", "title", "slug", "status", "body")
	builder.Values(
		n.Post.ID,
		n.Post.Title,
		n.Slug,
		n.Status,
		n.Post.Body,
	)
	query, args := builder.Build()

	_, err = tx.Exec(query, args...)
	if err != nil {
		if possibleErr, ok := err.(sqlite3.Error); ok {
			if possibleErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return bareknews.ErrDataAlreadyExist
			}
		}
		return errors.Wrap(err, "storage.news.save")
	}

	err = s.insertTags(tx, n)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "storage.news.save")
	}

	return nil
}

func (s News) GetById(id uuid.UUID) (*news.News, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug")
	builder.From("news")
	builder.Where(builder.Equal("id", id))

	query, args := builder.Build()
	row := s.Conn.QueryRow(query, args...)
	post := domain.Post{}
	slug := new(domain.Slug)
	status := new(domain.Status)
	err := row.Scan(&post.ID, &post.Title, status, &post.Body, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &news.News{}, sql.ErrNoRows
		} else {
			return &news.News{}, errors.Wrap(err, "storage.news.getById")
		}
	}

	tagsResult, err := s.getTags(post.ID)
	if err != nil {
		return &news.News{}, err
	}

	result := &news.News{
		Post:   post,
		Slug:   *slug,
		Status: *status,
		TagsID: tagsResult,
	}
	return result, nil
}

func (s News) Update(n news.News) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "storage.news.update")
	}

	defer tx.Rollback()

	builder := sqlbuilder.NewUpdateBuilder()
	builder.Update("news")
	builder.Set(
		builder.Assign("title", n.Post.Title),
		builder.Assign("body", n.Post.Body),
		builder.Assign("status", n.Status),
		builder.Assign("slug", n.Slug),
	)
	builder.Where(builder.Equal("id", n.Post.ID))
	query, args := builder.Build()
	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "storage.news.update")
	}

	// deleting relation between news and tags.
	err = s.deleteTags(tx, n.Post.ID)
	if err != nil {
		return err
	}

	// inserting new relation between news and tags if they
	// exist.
	err = s.insertTags(tx, n)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "storage.news.update")
	}

	return nil
}

func (s News) Delete(id uuid.UUID) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "storage.news.delete")
	}

	defer tx.Rollback()

	err = s.deleteTags(tx, id)
	if err != nil {
		return err
	}

	d := sqlbuilder.NewDeleteBuilder()
	d.DeleteFrom("news")
	d.Where(d.Equal("id", id))
	sql, args := d.Build()
	_, err = tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "storage.news.delete")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "storage.news.delete")
	}
	return nil
}

func (s News) Count(id uuid.UUID) (int, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select(builder.As("COUNT(id)", "c"))
	builder.From("news")
	builder.Where(builder.Equal("id", id))
	query, args := builder.Build()
	row := s.Conn.QueryRow(query, args...)

	var c int
	err := row.Scan(&c)
	if err != nil {
		return c, errors.Wrap(err, "storage.news.count")
	}

	if c == 0 {
		return c, sql.ErrNoRows
	}

	return c, nil
}

func (s News) GetAll() ([]news.News, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug")
	builder.From("news")
	query, args := builder.Build()
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "storage.news.getAll")
	}

	defer rows.Close()

	newsResults := make([]news.News, 0)

	for rows.Next() {
		post := domain.Post{}
		slug := new(domain.Slug)
		status := new(domain.Status)
		err = rows.Scan(&post.ID, &post.Title, status, &post.Body, slug)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "storage.news.getAll")
		}
		tagsResult, err := s.getTags(post.ID)
		if err != nil {
			return []news.News{}, err
		}
		newsResults = append(newsResults, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			TagsID: tagsResult,
		})
	}

	if rows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "storage.news.getAll")
	}

	return newsResults, nil
}

func (s News) GetAllByTopic(topic uuid.UUID) ([]news.News, error) {
	newsID, err := s.getNewsIds(topic)
	if err != nil {
		return []news.News{}, err
	}

	idNewsStr := make([]string, 0)

	for _, elem := range newsID {
		idNewsStr = append(idNewsStr, elem.String())
	}

	newsIdMark := sqlbuilder.List(idNewsStr)
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug")
	builder.From("news")
	builder.Where(builder.In("id", newsIdMark))
	newsQuery, newsArgs := builder.Build()

	newsRows, err := s.Conn.Query(newsQuery, newsArgs...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "storage.news.getAllByTopic")
	}

	defer newsRows.Close()

	newsResult := make([]news.News, 0)

	for newsRows.Next() {
		post := domain.Post{}
		slug := new(domain.Slug)
		status := new(domain.Status)
		err = newsRows.Scan(&post.ID, &post.Title, status, &post.Body, slug)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "storage.news.getAllByTopic")
		}
		tagsResult, err := s.getTags(post.ID)
		if err != nil {
			return []news.News{}, err
		}
		newsResult = append(newsResult, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			TagsID: tagsResult,
		})
	}

	if newsRows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "storage.news.getAllByTopic")
	}

	return newsResult, nil
}

func (s News) getNewsIds(tagsID uuid.UUID) ([]uuid.UUID, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Distinct()
	builder.Select("newsID")
	builder.From("news_tags")
	builder.Where(builder.Equal("tagsID", tagsID))
	query, args := builder.Build()
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "storage.news.getNewsIds")
	}

	defer rows.Close()

	newsID := make([]uuid.UUID, 0)

	for rows.Next() {
		id := uuid.UUID{}
		err = rows.Scan(&id)
		if err != nil {
			return []uuid.UUID{}, errors.Wrap(err, "storage.news.getNewsIds")
		}
		newsID = append(newsID, id)
	}

	if rows.Err() != nil {
		return []uuid.UUID{}, errors.Wrap(err, "storage.news.getNewsIds")
	}

	return newsID, nil
}

func (s News) insertTags(tx *sql.Tx, n news.News) error {
	for i := range n.TagsID {
		builder := sqlbuilder.NewInsertBuilder()
		builder.InsertInto("news_tags")
		builder.Cols("newsID", "tagsID")
		builder.Values(n.Post.ID, n.TagsID[i])
		query, args := builder.Build()

		_, err := tx.Exec(query, args...)
		if err != nil {
			if possibleErr, ok := err.(sqlite3.Error); ok {
				if possibleErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					return bareknews.ErrDataAlreadyExist
				}
			}
			return errors.Wrap(err, "storage.news.insertTagsReference")
		}
	}

	return nil
}

func (s News) deleteTags(tx *sql.Tx, id uuid.UUID) error {
	builderDel := sqlbuilder.NewDeleteBuilder()
	builderDel.DeleteFrom("news_tags")
	builderDel.Where(builderDel.Equal("newsID", id))
	query, args := builderDel.Build()

	_, err := tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "storage.news.deleteTagsReference")
	}

	return nil
}

func (s News) getTags(newsId uuid.UUID) ([]uuid.UUID, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Distinct()
	builder.Select("tagsID")
	builder.From("news_tags")
	builder.Where(builder.Equal("newsID", newsId))
	query, args := builder.Build()
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "storage.news.getTagsReference")
	}

	defer rows.Close()

	tagsResult := make([]uuid.UUID, 0)

	for rows.Next() {
		tagId := uuid.UUID{}
		err = rows.Scan(&tagId)
		if err != nil {
			return []uuid.UUID{}, errors.Wrap(err, "storage.news.getTagsReference")
		}
		tagsResult = append(tagsResult, tagId)
	}

	if rows.Err() != nil {
		return []uuid.UUID{}, errors.Wrap(err, "storage.news.getTags")
	}

	return tagsResult, nil
}
