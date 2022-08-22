package v1

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/Iiqbal2000/bareknews/news"
	newsdb "github.com/Iiqbal2000/bareknews/news/db"
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/Iiqbal2000/bareknews/tags"
	tagsdb "github.com/Iiqbal2000/bareknews/tags/db"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	DB       *sql.DB
}

func APIMux(config APIMuxConfig) http.Handler {
	app := web.NewApp(
		config.Shutdown,
		web.ContentTypeJSON(),
		web.CORS(),
		web.Errors(config.Log),
		web.Panics(),
	)

	// app.Mux.Get("/swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	// ))

	newsDB := newsdb.CreateStore(config.DB)
	tagsDB := tagsdb.CreateStore(config.DB)

	tagsSvc := tags.CreateSvc(tagsDB)
	newsSvc := news.CreateSvc(newsDB, tagsSvc)

	tagsHandler := tags.CreateHandler(tagsSvc, config.Log)
	newsHandler := news.CreateHandler(newsSvc, config.Log)

	app.Handle("POST", "/api/news", newsHandler.Create)
	app.Handle("GET", "/api/news", newsHandler.GetAll)
	app.Handle("GET", "/api/news/{newsId}", newsHandler.GetById)
	app.Handle("PUT", "/api/news/{newsId}", newsHandler.Update)
	app.Handle("DELETE", "/api/news/{newsId}", newsHandler.Delete)

	app.Handle("POST", "/api/tags", tagsHandler.Create)
	app.Handle("GET", "/api/tags", tagsHandler.GetAll)
	app.Handle("GET", "/api/tags/{tagId}", tagsHandler.GetById)
	app.Handle("PUT", "/api/tags/{tagId}", tagsHandler.Update)
	app.Handle("DELETE", "/api/tags/{tagId}", tagsHandler.Delete)
	return app
}
