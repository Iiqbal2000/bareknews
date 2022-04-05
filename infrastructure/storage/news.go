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

var post = sqlbuilder.NewStruct(new(domain.Post))

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
		n.Post.Status,
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
		builder.Assign("status", n.Post.Status),
		builder.Assign("slug", n.Slug),
	)
	builder.Where(builder.Equal("id", n.Post.ID))
	query, args := builder.Build()
	_, err = tx.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not found")
		} else {
			return fmt.Errorf("failure when updating a tag: %s", err.Error())
		}
	}

	err = s.deleteTags(tx, n.Post.ID)
	if err != nil {
		return err
	}

	err = s.insertTags(tx, n)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("faiure when commiting in updating news item")
	}

	return nil
}

func (s News) GetById(id string) (*news.News, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug")
	builder.From("news")
	builder.Where(builder.Equal("id", id))

	query, args := builder.Build()
	row := s.Conn.QueryRow(query, args...)
	post := domain.Post{}
	slug := new(domain.Slug)
	err := row.Scan(&post.ID, &post.Title, &post.Status, &post.Body, slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &news.News{}, fmt.Errorf("news is not found")
		} else {
			return &news.News{}, fmt.Errorf("failure when querying news: %s", err.Error())
		}
	}

	result := &news.News{Post: post, Slug: *slug}
	tagsResult, err := s.getTags(result.Post.ID)
	if err != nil {
		return &news.News{}, err
	}

	result.Tags = tagsResult
	return result, nil
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
		err = newsRows.Scan(&post.ID, &post.Title, &post.Status, &post.Body, slug)
		if err != nil {
			return []news.News{}, fmt.Errorf("(GetAll news operation) failure when querying news in iteration: %s", err.Error())
		}
		tagsResult, err := s.getTags(post.ID)
		if err != nil {
			return []news.News{}, err
		}
		newsResults = append(newsResults, news.News{
			Post: post,
			Slug: *slug,
			Tags: tagsResult,
		})
	}

	fmt.Println(newsResults)

	return newsResults, nil
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

func (s News) insertTags(tx *sql.Tx, n news.News) error {
	for i := range n.Tags {
		builder := sqlbuilder.NewInsertBuilder()
		builder.InsertInto("topics")
		builder.Cols("newsID", "tagsID")
		builder.Values(n.Post.ID, n.Tags[i].ID)
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
	builderDel.DeleteFrom("topics")
	builderDel.Where(builderDel.Equal("newsID", id))
	query, args := builderDel.Build()
	_, err := tx.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not found")
		} else {
			return fmt.Errorf("failure when deleting tags: %s", err.Error())
		}
	}

	return nil
}

func (s News) getTags(newsId uuid.UUID) ([]domain.Tags, error) {
	builder := sqlbuilder.NewSelectBuilder()
	builder.Distinct()
	builder.Select("tags.ID", "tags.name", "tags.slug")
	builder.From("news")
	builder.Join("topics", builder.Equal("topics.newsID", newsId))
	builder.Join("tags", "tags.ID = topics.tagsID")
	query, args := builder.Build()
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		return []domain.Tags{}, fmt.Errorf("failure when querying tags in news: %s", err.Error())
	}

	defer rows.Close()

	tagsResult := make([]domain.Tags, 0)

	for rows.Next() {
		t := domain.Tags{}
		err = rows.Scan(tags.Addr(&t)...)
		if err != nil {
			return []domain.Tags{}, fmt.Errorf("failure when querying tags in the news iteration: %s", err.Error())
		}
		tagsResult = append(tagsResult, t)
	}

	return tagsResult, nil
}
