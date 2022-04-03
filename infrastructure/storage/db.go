package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var tagsquery = `CREATE TABLE IF NOT EXISTS tags(
	ID VARCHAR (127) PRIMARY KEY UNIQUE,
	name VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL
)`

var newsQuery = `CREATE TABLE IF NOT EXISTS news(
	ID CHAR (127) PRIMARY KEY UNIQUE,
	title VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL UNIQUE,
	status VARCHAR (127) NOT NULL,
	body TEXT NOT NULL
)`

var topicsquery = `CREATE TABLE IF NOT EXISTS topics(
	newsID VARCHAR (127) NOT NULL,
	tagsID VARCHAR (127) NOT NULL,
	FOREIGN KEY(newsID) REFERENCES news(id) ON DELETE CASCADE,
	FOREIGN KEY(tagsID) REFERENCES tags(id) ON DELETE CASCADE
)`

func Run(dbpath string, dropTable bool) *sql.DB {
	if _, err := os.Stat(dbpath); err != nil {
		log.Fatal("The DB file is not found")
	}

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal("failure when opening db connection: ", err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("failure when starting transaction: ", err.Error())
	}

	defer tx.Rollback()

	if dropTable {
		_, err = tx.Exec("DROP TABLE IF EXISTS tags;")
		if err != nil {
			log.Fatal("failure when drop tags table: ", err.Error())
		}

		_, err = tx.Exec("DROP TABLE IF EXISTS news;")
		if err != nil {
			log.Fatal("failure when drop news table: ", err.Error())
		}

		_, err = tx.Exec("DROP TABLE IF EXISTS topics;")
		if err != nil {
			log.Fatal("failure when drop topics table: ", err.Error())
		}
	}

	_, err = tx.Exec(tagsquery)
	if err != nil {
		log.Fatal("failure when creating tags table: ", err.Error())
	}

	_, err = tx.Exec(newsQuery)
	if err != nil {
		log.Fatal("failure when creating news table: ", err.Error())
	}

	_, err = tx.Exec(topicsquery)
	if err != nil {
		log.Fatal("failure when creating topics table: ", err.Error())
	}

	if err = tx.Commit(); err != nil {
		log.Fatal("faiure commiting")
	}
	log.Println("successfully connected to database")
	return db
}
