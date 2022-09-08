package db

import (
	"context"
	"database/sql"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/Iiqbal2000/bareknews/news/db")

type Store struct {
	conn *sql.DB
}

func CreateStore(conn *sql.DB) Store {
	return Store{conn: conn}
}

func (s Store) Save(ctx context.Context, n news.News) error {
	ctx, span := tracer.Start(ctx, "news.db.Save")
	defer span.End()

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer tx.Rollback()

	builder := sqlbuilder.NewInsertBuilder()
	builder.InsertInto("news")
	builder.Cols(
		"id",
		"title",
		"slug",
		"status",
		"body",
		"date_created",
		"date_updated",
	)
	builder.Values(
		n.Post.ID,
		n.Post.Title,
		n.Slug,
		n.Status,
		n.Post.Body,
		n.DateCreated,
		n.DateUpdated,
	)
	query, args := builder.Build()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if possibleErr, ok := err.(sqlite3.Error); ok {
			if possibleErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return bareknews.ErrDataAlreadyExist
			}
		}
		return errors.Wrap(err, "exec the query")
	}

	err = s.insertNewsTagsRelation(ctx, tx, n)
	if err != nil {
		return errors.Wrap(err, "could not insert news-tags relation")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (s Store) GetById(ctx context.Context, id uuid.UUID) (*news.News, error) {
	ctx, span := tracer.Start(ctx, "news.db.GetById")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()
	builder.Select("id", "title", "status", "body", "slug", "date_created", "date_updated")
	builder.From("news")
	builder.Where(builder.Equal("id", id))

	query, args := builder.Build()
	row := s.conn.QueryRowContext(ctx, query, args...)
	post := bareknews.Post{}
	slug := new(bareknews.Slug)
	status := new(bareknews.Status)
	dateCreated := new(int64)
	updateCreated := new(int64)

	err := row.Scan(&post.ID, &post.Title, status, &post.Body, slug, dateCreated, updateCreated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &news.News{}, sql.ErrNoRows
		} else {
			return &news.News{}, errors.Wrap(err, "scan a news item")
		}
	}

	tagIdResults, err := s.getAllTagIds(ctx, post.ID)
	if err != nil {
		return &news.News{}, errors.Wrap(err, "could not get tag ids")
	}

	result := &news.News{
		Post:   post,
		Slug:   *slug,
		Status: *status,
		TagsID: tagIdResults,
		DateCreated: *dateCreated,
		DateUpdated: *updateCreated,
	}
	return result, nil
}

func (s Store) Update(ctx context.Context, n news.News) error {
	ctx, span := tracer.Start(ctx, "news.db.Update")
	defer span.End()

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer tx.Rollback()

	builder := sqlbuilder.NewUpdateBuilder()
	builder.Update("news")
	builder.Set(
		builder.Assign("title", n.Post.Title),
		builder.Assign("body", n.Post.Body),
		builder.Assign("status", n.Status),
		builder.Assign("slug", n.Slug),
		builder.Assign("date_updated", n.DateUpdated),
	)

	builder.Where(builder.Equal("id", n.Post.ID))

	query, args := builder.Build()
	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "exec the query")
	}

	// deleting relation between news and tags.
	err = s.deleteNewsTagsRelation(ctx, tx, n.Post.ID)
	if err != nil {
		return errors.Wrap(err, "could not delete news-tags relation")
	}

	// inserting new relation between news and tags if they
	// exist.
	err = s.insertNewsTagsRelation(ctx, tx, n)
	if err != nil {
		return errors.Wrap(err, "could not insert news-tags relation")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (s Store) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "news.db.Delete")
	defer span.End()

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer tx.Rollback()

	err = s.deleteNewsTagsRelation(ctx, tx, id)
	if err != nil {
		return errors.Wrap(err, "could not delete news-tags relation")
	}

	d := sqlbuilder.NewDeleteBuilder()
	d.DeleteFrom("news")
	d.Where(d.Equal("id", id))
	sql, args := d.Build()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "exec the query")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}
	return nil
}

func (s Store) Count(ctx context.Context, id uuid.UUID) (int, error) {
	ctx, span := tracer.Start(ctx, "news.db.Count")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()

	builder.Select(builder.As("COUNT(id)", "c"))
	builder.From("news")
	builder.Where(builder.Equal("id", id))

	query, args := builder.Build()
	row := s.conn.QueryRowContext(ctx, query, args...)

	var result int
	err := row.Scan(&result)
	if err != nil {
		return result, errors.Wrap(err, "scan a total amount of news")
	}

	if result == 0 {
		return result, sql.ErrNoRows
	}

	return result, nil
}

func (s Store) GetAll(ctx context.Context, cursor int64, limit int) ([]news.News, error) {
	ctx, span := tracer.Start(ctx, "news.db.GetAll")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("id", "title", "status", "body", "slug", "date_created", "date_updated")
	builder.From("news")
	if cursor != 0 {
		builder.Where(builder.LessThan("date_created", cursor))
	}

	builder.OrderBy("date_created").Desc()

	if limit == 0 {
		limit = 2
	}

	builder.Limit(limit)

	query, args := builder.Build()

	rows, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "exec the query")
	}

	defer rows.Close()

	newsResults := make([]news.News, 0)
	postIds := make([]uuid.UUID, 0)

	for rows.Next() {
		post := bareknews.Post{}
		slug := new(bareknews.Slug)
		status := new(bareknews.Status)
		dateCreated := new(int64)
		dateUpdted := new(int64)

		err = rows.Scan(
			&post.ID,
			&post.Title,
			status,
			&post.Body,
			slug,
			dateCreated,
			dateUpdted,
		)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "scan a news item")
		}

		postIds = append(postIds, post.ID)

		newsResults = append(newsResults, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			DateCreated: *dateCreated,
			DateUpdated: *dateUpdted,
		})
	}

	if rows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "failed get items during iteration")
	}

	tagIdBucket, err := s.getAllNewsTagsIds(ctx, postIds)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "could not get tag ids")
	}

	// Reconstruct news with addition TagsID property.
	for i, elem := range newsResults {
		newsResults[i] = news.News{
			Post:   elem.Post,
			Status: elem.Status,
			Slug:   elem.Slug,
			TagsID: tagIdBucket[elem.Post.ID],
			DateCreated: elem.DateCreated,
			DateUpdated: elem.DateUpdated,
		}
	}

	return newsResults, nil
}

func (s Store) GetAllByPagination(ctx context.Context, cursor int64, limit int) ([]news.News, error) {
	ctx, span := tracer.Start(ctx, "news.db.GetAllByPagination")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()

	builder.Select(
		"id",
		"title",
		"status",
		"body",
		"slug",
		"date_created",
		"date_updated",
	)
	builder.From("news")

	if cursor != 0 {
		builder.Where(builder.LessThan("date_created", cursor))
	}

	builder.OrderBy("date_created").Desc()

	if limit == 0 {
		limit = 2
	}

	builder.Limit(limit)

	query, args := builder.Build()

	rows, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "exec the query")
	}

	defer rows.Close()

	newsResults := make([]news.News, 0)
	postIds := make([]uuid.UUID, 0)

	for rows.Next() {
		post := bareknews.Post{}
		slug := new(bareknews.Slug)
		status := new(bareknews.Status)
		dateCreated := new(int64)
		dateUpdted := new(int64)

		err = rows.Scan(
			&post.ID,
			&post.Title,
			status,
			&post.Body,
			slug,
			dateCreated,
			dateUpdted,
		)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "scan a news item")
		}

		postIds = append(postIds, post.ID)

		newsResults = append(newsResults, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			DateCreated: *dateCreated,
			DateUpdated: *dateUpdted,
		})
	}

	if rows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "failed get items during iteration")
	}

	tagIdBucket, err := s.getAllNewsTagsIds(ctx, postIds)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "could not get tag ids")
	}

	// Reconstruct news with addition TagsID property.
	for i, elem := range newsResults {
		newsResults[i] = news.News{
			Post:   elem.Post,
			Status: elem.Status,
			Slug:   elem.Slug,
			TagsID: tagIdBucket[elem.Post.ID],
			DateCreated: elem.DateCreated,
			DateUpdated: elem.DateUpdated,
		}
	}

	return newsResults, nil
}

func (s Store) GetAllByTopic(ctx context.Context, topic uuid.UUID, cursor int64, limit int) ([]news.News, error) {
	ctx, span := tracer.Start(ctx, "news.db.GetAllByTopic")
	defer span.End()

	newsIDs, err := s.getAllNewsIds(ctx, topic)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "could not get news ids")
	}

	newsIdsStr := make([]string, 0)

	for _, elem := range newsIDs {
		newsIdsStr = append(newsIdsStr, elem.String())
	}

	newsIdMark := sqlbuilder.List(newsIdsStr)
	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("id", "title", "status", "body", "slug","date_created", "date_updated",)
	builder.From("news")
	builder.Where(builder.In("id", newsIdMark))
	if cursor != 0 {
		builder.Where(builder.LessThan("date_created", cursor))
	}

	builder.OrderBy("date_created").Desc()

	if limit == 0 {
		limit = 2
	}

	builder.Limit(limit)

	newsQuery, newsArgs := builder.Build()

	newsRows, err := s.conn.QueryContext(ctx, newsQuery, newsArgs...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "exec the query")
	}

	defer newsRows.Close()

	newsResult := make([]news.News, 0)
	postIds := make([]uuid.UUID, 0)

	for newsRows.Next() {
		post := bareknews.Post{}
		slug := new(bareknews.Slug)
		status := new(bareknews.Status)
		dateCreated := new(int64)
		dateUpdted := new(int64)

		err = newsRows.Scan(
			&post.ID,
			&post.Title,
			status,
			&post.Body,
			slug,
			dateCreated,
			dateUpdted,
		)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "scan a news item")
		}

		postIds = append(postIds, post.ID)

		newsResult = append(newsResult, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			DateCreated: *dateCreated,
			DateUpdated: *dateUpdted,
		})
	}

	if newsRows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "failed get items during iteration")
	}

	tagIdsBucket, err := s.getAllNewsTagsIds(ctx, postIds)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "could not get tag ids in bucket")
	}

	// Reconstruct news with addition TagsID property.
	for i, elem := range newsResult {
		newsResult[i] = news.News{
			Post:   elem.Post,
			Status: elem.Status,
			Slug:   elem.Slug,
			TagsID: tagIdsBucket[elem.Post.ID],
			DateCreated: elem.DateCreated,
			DateUpdated: elem.DateUpdated,
		}
	}

	return newsResult, nil
}

func (s Store) GetAllByStatus(ctx context.Context, status bareknews.Status, cursor int64, limit int) ([]news.News, error) {
	ctx, span := tracer.Start(ctx, "news.db.GetAllByStatus")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("id", "title", "status", "body", "slug", "date_created", "date_updated")
	builder.From("news")
	builder.Where(builder.Equal("status", status))
	if cursor != 0 {
		builder.Where(builder.LessThan("date_created", cursor))
	}

	builder.OrderBy("date_created").Desc()

	if limit == 0 {
		limit = 2
	}

	builder.Limit(limit)

	newsQuery, newsArgs := builder.Build()

	newsRows, err := s.conn.QueryContext(ctx, newsQuery, newsArgs...)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "exec the query")
	}

	defer newsRows.Close()

	newsResult := make([]news.News, 0)
	postIds := make([]uuid.UUID, 0)

	for newsRows.Next() {
		post := bareknews.Post{}
		slug := new(bareknews.Slug)
		status := new(bareknews.Status)
		dateCreated := new(int64)
		dateUpdted := new(int64)

		err = newsRows.Scan(
			&post.ID,
			&post.Title,
			status,
			&post.Body,
			slug,
			dateCreated,
			dateUpdted,
		)
		if err != nil {
			return []news.News{}, errors.Wrap(err, "scan a news item")
		}

		postIds = append(postIds, post.ID)

		newsResult = append(newsResult, news.News{
			Post:   post,
			Status: *status,
			Slug:   *slug,
			DateCreated: *dateCreated,
			DateUpdated: *dateUpdted,
		})
	}

	if newsRows.Err() != nil {
		return []news.News{}, errors.Wrap(err, "failed get items during iteration")
	}

	tagIdsBucket, err := s.getAllNewsTagsIds(ctx, postIds)
	if err != nil {
		return []news.News{}, errors.Wrap(err, "could not get tag ids in bucket")
	}

	// Reconstruct news with addition TagsID property.
	for i, elem := range newsResult {
		newsResult[i] = news.News{
			Post:   elem.Post,
			Status: elem.Status,
			Slug:   elem.Slug,
			TagsID: tagIdsBucket[elem.Post.ID],
			DateCreated: elem.DateCreated,
			DateUpdated: elem.DateUpdated,
		}
	}

	return newsResult, nil
}

func (s Store) getAllNewsIds(ctx context.Context, tagsID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "news.db.getAllNewsIds")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()
	builder.Distinct()
	builder.Select("newsID")
	builder.From("news_tags")
	builder.Where(builder.Equal("tagsID", tagsID))
	query, args := builder.Build()

	rows, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "exec the query")
	}

	defer rows.Close()

	newsIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		id := uuid.UUID{}
		err = rows.Scan(&id)
		if err != nil {
			return []uuid.UUID{}, errors.Wrap(err, "scan a news id")
		}
		newsIDs = append(newsIDs, id)
	}

	if rows.Err() != nil {
		return []uuid.UUID{}, errors.Wrap(err, "failed get items during iteration")
	}

	return newsIDs, nil
}

func (s Store) insertNewsTagsRelation(ctx context.Context, tx *sql.Tx, nws news.News) error {
	_, span := tracer.Start(ctx, "news.db.insertNewsTagsRelation")
	defer span.End()

	// to avoid incomplete input error
	if len(nws.TagsID) == 0 {
		return nil
	}

	builder := sqlbuilder.NewInsertBuilder()

	builder.InsertInto("news_tags")
	builder.Cols("newsID", "tagsID")

	for i := range nws.TagsID {
		builder.Values(nws.Post.ID, nws.TagsID[i])
	}

	query, args := builder.Build()

	_, err := tx.Exec(query, args...)
	if err != nil {
		if possibleErr, ok := err.(sqlite3.Error); ok {
			if possibleErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return bareknews.ErrDataAlreadyExist
			}
		}
		return errors.Wrap(err, "exec the query")
	}
	return nil
}

func (s Store) deleteNewsTagsRelation(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	_, span := tracer.Start(ctx, "news.db.deleteNewsTagsRelation")
	defer span.End()

	builderDel := sqlbuilder.NewDeleteBuilder()

	builderDel.DeleteFrom("news_tags")
	builderDel.Where(builderDel.Equal("newsID", id))
	query, args := builderDel.Build()

	_, err := tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "exec the query")
	}

	return nil
}

func (s Store) getAllNewsTagsIds(ctx context.Context, newsIds []uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	_, span := tracer.Start(ctx, "news.db.getAllNewsTagsIds")
	defer span.End()
	
	idstr := make([]string, 0)

	for _, elem := range newsIds {
		idstr = append(idstr, elem.String())
	}

	builder := sqlbuilder.NewSelectBuilder()

	l := sqlbuilder.List(idstr)
	builder.Select("newsID", "tagsID")
	builder.From("news_tags")
	builder.Where(builder.In("newsID", l))
	query, args := builder.Build()

	rows, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return make(map[uuid.UUID][]uuid.UUID), errors.Wrap(err, "exec the query")
	}

	defer rows.Close()

	itemBatch := make(map[uuid.UUID][]uuid.UUID)

	for rows.Next() {
		tagId := uuid.UUID{}
		newsId := uuid.UUID{}

		err = rows.Scan(&newsId, &tagId)
		if err != nil {
			return make(map[uuid.UUID][]uuid.UUID), errors.Wrap(err, "scan news and tag id")
		}

		itemBatch[newsId] = append(itemBatch[newsId], tagId)
	}

	if rows.Err() != nil {
		return make(map[uuid.UUID][]uuid.UUID), errors.Wrap(err, "failed get items during iteration")
	}

	return itemBatch, nil
}

func (s Store) getAllTagIds(ctx context.Context, newsId uuid.UUID) ([]uuid.UUID, error) {
	_, span := tracer.Start(ctx, "news.db.getAllTagIds")
	defer span.End()

	builder := sqlbuilder.NewSelectBuilder()

	builder.Distinct()
	builder.Select("tagsID")
	builder.From("news_tags")
	builder.Where(builder.Equal("newsID", newsId))
	query, args := builder.Build()

	rows, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "exec the query")
	}

	defer rows.Close()

	tagsResult := make([]uuid.UUID, 0)

	for rows.Next() {
		tagId := uuid.UUID{}
		err = rows.Scan(&tagId)
		if err != nil {
			return []uuid.UUID{}, errors.Wrap(err, "scan a tag id")
		}
		tagsResult = append(tagsResult, tagId)
	}

	if rows.Err() != nil {
		return []uuid.UUID{}, errors.Wrap(err, "failed get items during iteration")
	}

	return tagsResult, nil
}
