package main

import (
	"fmt"
	"log"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/infrastructure/storage"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

const dbfile = "./bareknews.db"

func main() {
	dbConn := storage.Run(dbfile, false)
	builder := sqlbuilder.NewInsertBuilder()
	builder.InsertInto("news")
	builder.Cols("id", "title", "slug", "status", "body")
	builder.Values(
		uuid.New(),
		uuid.New().String(),
		uuid.New().String(),
		"draft",
		"nsskoskoskos",
	)
	query, args := builder.Build()

	_, err := dbConn.Exec(query, args...)
	if err != nil {
		log.Fatal("failure when inserting ", err.Error())
	}
	newsRows, err := dbConn.Query("select id, title, status, body, slug from news")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer newsRows.Close()
	post := make([]domain.Posts, 0)

	for newsRows.Next() {
		p := domain.Posts{}
		var slug domain.Slug
		err = newsRows.Scan(&p.ID, &p.Title, &p.Status, &p.Body, &slug)
		if err != nil {
			log.Fatal("failure scanning: ", err.Error())
		}

		fmt.Println("post ", p)

		post = append(post, p)

		if err != nil {
			log.Println(err.Error())
		}
	}

	if newsRows.Err() != nil {
		log.Fatal(newsRows.Err().Error())
	}

	// newsDB := storage.News{Conn: dbConn}
	// tagDB := storage.Tag{Conn: dbConn}

	// taggingSvc := tagging.New(tagDB)
	// postingSvc := posting.New(newsDB, taggingSvc)
}
