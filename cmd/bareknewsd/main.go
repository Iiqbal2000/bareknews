package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/Iiqbal2000/bareknews/docs"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/Iiqbal2000/bareknews/pkg/sqlite3"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

const dbfile = "./bareknews.db"
const address = ":3333"

// @title Bareknews API
// @version 1.0
// @description This is a sample server Bareknews server.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3333
// @BasePath /api
func main() {
	r := chi.NewRouter()

	dbConn := sqlite3.Run(dbfile, true)
	newsDB := sqlite3.News{Conn: dbConn}
	tagDB := sqlite3.Tag{Conn: dbConn}

	taggingSvc := tags.CreateSvc(tagDB)
	postingSvc := news.CreateSvc(newsDB, taggingSvc)

	tagHandler := tags.Restapi{Service: taggingSvc}
	newsHandler := news.Restapi{Service: postingSvc}

	r.Use(ContentTypeJSON)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))
	
	r.Route("/api", func(r chi.Router) {
		r.Route("/tags", tagHandler.Route)
		r.Route("/news", newsHandler.Route)
	})

	s := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}

	log.Printf("Starting server at PORT %s\n", address)

	s.ListenAndServe()
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}
