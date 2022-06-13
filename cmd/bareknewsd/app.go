package main

import (
	"log"

	"github.com/Iiqbal2000/bareknews/news"
	"github.com/Iiqbal2000/bareknews/pkg/sqlite3"
	"github.com/Iiqbal2000/bareknews/tags"
)

const dbfile = "./bareknews.db"

type App struct {
	TaggingSvc tags.Service
	PostingSvc news.Service
}

func RunApp() App {
	dbConn := sqlite3.Run(dbfile, true)
	newsDB := sqlite3.News{Conn: dbConn}
	tagDB := sqlite3.Tag{Conn: dbConn}

	taggingSvc := tags.CreateSvc(tagDB)
	postingSvc := news.CreateSvc(newsDB, taggingSvc)

	log.Println("Starting APP")

	return App{
		TaggingSvc: taggingSvc,
		PostingSvc: postingSvc,
	}
}
