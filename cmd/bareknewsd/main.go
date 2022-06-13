package main

import (
	"net/http"

	_ "github.com/Iiqbal2000/bareknews/docs"
	"github.com/Iiqbal2000/bareknews/news"
	"github.com/Iiqbal2000/bareknews/pkg/restapi"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

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

	app := RunApp()

	tagHandler := tags.Restapi{Service: app.TaggingSvc}
	newsHandler := news.Restapi{Service: app.PostingSvc}

	r.Use(ContentTypeJSON)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Route("/tags", tagHandler.Route)
		r.Route("/news", newsHandler.Route)
	})

	restapi.CreateAndRunServer(r, address)
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}
