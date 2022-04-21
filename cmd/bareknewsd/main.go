package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Iiqbal2000/bareknews/infrastructure/handler"
	"github.com/Iiqbal2000/bareknews/infrastructure/storage"
	"github.com/Iiqbal2000/bareknews/services/posting"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/go-chi/chi/v5"
)

const dbfile = "./bareknews.db"

func main() {
	r := chi.NewRouter()

	dbConn := storage.Run(dbfile, true)
	newsDB := storage.News{Conn: dbConn}
	tagDB := storage.Tag{Conn: dbConn}

	taggingSvc := tagging.New(tagDB)
	postingSvc := posting.New(newsDB, taggingSvc)

	tagHandler := handler.Tags{Service: taggingSvc}
	newsHandler := handler.News{Service: postingSvc}

	r.Use(ContentTypeJSON)

	r.Route("/tags", tagHandler.Route)
	r.Route("/news", newsHandler.Route)

	address := ":3333"

	s := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}

	log.Printf("Starting server at %s\n", address)

	s.ListenAndServe()
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}
