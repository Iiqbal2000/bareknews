package sqlite3

import (
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed schema/*.sql
var schema embed.FS

//go:embed sql/seed.sql
var seedDoc string

// Config
type Config struct {
	URI            string
	Log            *zap.SugaredLogger
	DropTableFirst bool
}

// Run
func Run(c Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", c.URI)
	if err != nil {
		return nil, errors.Wrap(err, "failure when opening db connection")
	}

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(schema)

	if c.DropTableFirst {
		if err := goose.Reset(db, "schema"); err != nil {
			return nil, errors.Wrap(err, "could not perform resetting the migration")
		}
	}

	if err := goose.Up(db, "schema"); err != nil {
		return nil, errors.Wrap(err, "could not perform schema migration")
	}

	if err := goose.Version(db, "schema"); err != nil {
		return nil, errors.Wrap(err, "could not prints the current version of the db migration")
	}

	if c.Log != nil {
		c.Log.Infow("startup", "status", "successfully migrate the db schema")
	}

	return db, nil
}

// RunForTesting
func RunForTesting(URI string, Log *zap.SugaredLogger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", URI)
	if err != nil {
		return nil, errors.Wrap(err, "failure when opening db connection")
	}

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(schema)

	if err := goose.Reset(db, "schema"); err != nil {
		return nil, errors.Wrap(err, "could not perform resetting the migration")
	}

	if err := goose.Up(db, "schema"); err != nil {
		return nil, errors.Wrap(err, "could not perform schema migration")
	}

	if err := goose.Version(db, "schema"); err != nil {
		return nil, errors.Wrap(err, "could not prints the current version of the db migration")
	}

	Log.Info("Database is ready for test")

	return db, nil
}

// Seed
func Seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seedDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
