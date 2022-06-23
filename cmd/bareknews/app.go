package main

import (
	"net/http"

	"github.com/Iiqbal2000/bareknews/news"
	newsdb "github.com/Iiqbal2000/bareknews/news/db"
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/Iiqbal2000/bareknews/tags"
	tagsdb "github.com/Iiqbal2000/bareknews/tags/db"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func APIMux(cfg web.APIMuxConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(web.ContentTypeJSON)

	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))

	newsDB := newsdb.CreateStore(cfg.DB)
	tagDB := tagsdb.CreateStore(cfg.DB)

	tr := tags.Route(cfg, tagDB)
	nr := news.Route(cfg, newsDB, tr.Svc)

	mux.Route("/api", func(r chi.Router) {
		r.Route("/tags", tr.Router)
		r.Route("/news", nr)
	})

	return mux
}
