package it

import (
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews/pkg/logger"
	"github.com/Iiqbal2000/bareknews/pkg/sqlite3"
	"go.uber.org/zap"
)

// Test
type Test struct {
	Log      *zap.SugaredLogger
	DB       *sql.DB
	Teardown func()
	t        *testing.T
}

// RunDepedencies
func RunDepedencies(t *testing.T) Test {
	log, err := logger.New("e2e-test")
	if err != nil {
		t.Fatalf("error when initializing log: %s\n", err.Error())
		return Test{}
	}

	db, err := sqlite3.RunForTesting("./bareknews-integration-test.db", log)
	if err != nil {
		t.Fatalf("error when initializing : %s\n", err.Error())
		return Test{}
	}

	if err = sqlite3.Seed(db); err != nil {
		t.Fatalf("error when seeding the database : %s\n", err.Error())
		return Test{}
	}

	teardown := func() {
		t.Helper()
		db.Close()

		log.Sync()
	}

	return Test{
		Log:      log,
		DB:       db,
		Teardown: teardown,
		t:        t,
	}
}
