package web

import (
	"database/sql"

	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Log *zap.SugaredLogger
	DB  *sql.DB
}
