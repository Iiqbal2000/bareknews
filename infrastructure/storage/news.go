package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type News struct {
	Conn *sql.DB
}

func (s News) Save(n news.News) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failure when starting transaction: %s", err.Error())
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
		return fmt.Errorf("failure when inserting a news: %s", err.Error())
	}

	err = s.insertTags(tx, n)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("faiure commiting in news table")
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
			return &news.News{}, fmt.Errorf("failure when querying news: %s", err.Error())
		}
	}
	
	tagsResult, err := s.getTags(post.ID)
	if err != nil {
		return &news.News{}, err
	}

	result := &news.News{
		Post: post, 
		Slug: *slug, 
		Status: *status,
		TagsID: tagsResult,
	}
	return result, nil
}

func (s News) Update(n news.News) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failure when starting transaction in updating a news: %s", err.Error())
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
		return fmt.Errorf("failure when updating a tag: %s", err.Error())
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
		return fmt.Errorf("faiure when commit in updating news item")
	}

	return nil
}

func (s News) Delete(id uuid.UUID) error {
	tx, err := s.Conn.Begin()
	if err != nil {
		return fmt.Errorf("(news deleting operation) failure when starting the transaction: %s", err.Error())
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
		return fmt.Errorf("(news deleting operation) failure when deleting a news: %s", err.Error())
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("(news deleting operation) faiure when committing in updating news item")
	}
	return nil
}

func (s News) GetAll() ([]news.News, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug")
	builder.From("news")
	query, args := builder.Build()
	newsRows, err := s.Conn.Query(query, args...)
	if err != nil {
		return []news.News{}, fmt.Errorf("(GetAll news operation) failure when querying news: %s", err.Error())
	}

	defer newsRows.Close()

	newsResults := make([]news.News, 0)

	for newsRows.Next() {
		post := domain.Post{}
		slug := new(domain.Slug)
		status := new(domain.Status)
		err = newsRows.Scan(&post.ID, &post.Title, status, &post.Body, slug)
		if err != nil {
			return []news.News{}, fmt.Errorf("(GetAll news operation) failure when querying news in iteration: %s", err.Error())
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

	return newsResults, nil
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
			return fmt.Errorf("failure when inserting tags: %s", err.Error())
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
		return fmt.Errorf("failure when deleting tags: %s", err.Error())
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
		return []uuid.UUID{}, fmt.Errorf("failure when querying tags in news: %s", err.Error())
	}

	defer rows.Close()

	tagsResult := make([]uuid.UUID, 0)

	for rows.Next() {
		tagId := uuid.UUID{}
		err = rows.Scan(&tagId)
		if err != nil {
			return []uuid.UUID{}, fmt.Errorf("failure when querying tags in the news iteration: %s", err.Error())
		}
		tagsResult = append(tagsResult, tagId)
	}

	return tagsResult, nil
}
