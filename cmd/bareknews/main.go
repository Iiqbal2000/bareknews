package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Iiqbal2000/bareknews/docs"
	"github.com/Iiqbal2000/bareknews/pkg/logger"
	"github.com/Iiqbal2000/bareknews/pkg/sqlite3"
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

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
	// Construct the application logger.
	log, err := logger.New("BAREKNEWS-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := struct {
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3333"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
		DB string `conf:"default:./bareknews.db"`
	}{
		// DB: dbfile,
	}

	const prefix = "NEWS"
	_, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service")
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Infow("startup", "config", out)

	dbConn, err := sqlite3.Run(sqlite3.Config{
		URI: cfg.DB,
	}, true)
	if err != nil {
		return errors.Wrap(err, "failed to connect db")
	}

	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB)
		dbConn.Close()
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing web API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiMux := APIMux(web.APIMuxConfig{
		Log: log,
		DB:  dbConn,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
