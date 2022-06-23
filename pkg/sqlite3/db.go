package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const tagsquery = `CREATE TABLE IF NOT EXISTS tags(
	ID VARCHAR (127) PRIMARY KEY UNIQUE,
	name VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL UNIQUE
)`

const newsQuery = `CREATE TABLE IF NOT EXISTS news(
	ID CHAR (127) PRIMARY KEY UNIQUE,
	title VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL UNIQUE,
	status VARCHAR (127) NOT NULL,
	body TEXT NOT NULL
)`

const news_tagsquery = `CREATE TABLE IF NOT EXISTS news_tags(
	newsID VARCHAR (127) NOT NULL,
	tagsID VARCHAR (127) NOT NULL,
	FOREIGN KEY(newsID) REFERENCES news(id) ON DELETE CASCADE,
	FOREIGN KEY(tagsID) REFERENCES tags(id) ON DELETE CASCADE
)`

type Config struct {
	URI string
}

func Run(c Config, dropTable bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", c.URI)
	if err != nil {
		return nil, errors.Wrap(err, "failure when opening db connection")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failure when starting transaction")
	}

	defer tx.Rollback()

	if dropTable {
		_, err = tx.Exec("DROP TABLE IF EXISTS tags;")
		if err != nil {
			return nil, errors.Wrap(err, "failure when drop tags table")
		}

		_, err = tx.Exec("DROP TABLE IF EXISTS news;")
		if err != nil {
			return nil, errors.Wrap(err, "failure when drop news table")
		}

		_, err = tx.Exec("DROP TABLE IF EXISTS news_tags;")
		if err != nil {
			return nil, errors.Wrap(err, "failure when drop news_tags table")
		}
	}

	_, err = tx.Exec(tagsquery)
	if err != nil {
		return nil, errors.Wrap(err, "failure when creating tags table")
	}

	_, err = tx.Exec(newsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failure when creating news table")
	}

	_, err = tx.Exec(news_tagsquery)
	if err != nil {
		return nil, errors.Wrap(err, "failure when creating news_tags table")
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "faiure when commiting the queries")
	}
	return db, nil
}
