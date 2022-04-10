package main

import (
	"net/http"

	"github.com/Iiqbal2000/bareknews/infrastructure/handler"
	"github.com/Iiqbal2000/bareknews/infrastructure/storage"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/go-chi/chi/v5"
)

const dbfile = "./bareknews.db"

func main() {
	r := chi.NewRouter()

	dbConn := storage.Run(dbfile, true)
	// newsDB := storage.News{Conn: dbConn}
	tagDB := storage.Tag{Conn: dbConn}

	taggingSvc := tagging.New(tagDB)
	// postingSvc := posting.New(newsDB, taggingSvc)

	tagHandler := handler.Tags{Service: taggingSvc}
	
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			h.ServeHTTP(w, r)
		})
	})

	r.Route("/tags", tagHandler.Route)

	http.ListenAndServe(":3333", r)
}
