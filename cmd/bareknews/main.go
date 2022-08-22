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
	v1 "github.com/Iiqbal2000/bareknews/pkg/v1"
	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

const prefix = "NEWS"
const appName = "BAREKNEWS-API"
const tracerProviderUrl = "http://localhost:14268/api/traces"

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
	log, err := logger.New(appName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err.Error())
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
	}{}

	_, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// Starting The Supports

	// Starting a tracer provider support.
	tp, err := NewJaegerTracerProvider(tracerProviderUrl)
	if err != nil {
		return errors.Wrap(err, "starting the tracer provider")
	}

	defer func() {
		log.Info("shutdown the tracer provider")
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Errorf("Error shutting down tracer provider: %v", err)
		}
	}()

	log.Info("starting the app")

	// Generating the config.
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	log.Infow("config of app", "config", out)

	// Starting a database support.
	dbConn, err := sqlite3.Run(sqlite3.Config{
		URI:            cfg.DB,
		DropTableFirst: true,
		Log:            log,
	})

	if err != nil {
		return errors.Wrap(err, "failed to connect db")
	}

	defer func() {
		log.Infow("shutdown the database", "host", cfg.DB)
		if err := dbConn.Close(); err != nil {
			log.Errorf("Error shutting down database: %v", err)
		}
	}()

	// =========================================================================
	// Start API Service

	log.Info("initializing web API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	muxConfig := v1.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		DB:       dbConn,
	}

	mux := v1.APIMux(muxConfig)

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      otelhttp.NewHandler(app, "request"),
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
		log.Infow("api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown Handler
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown: // Graceful shutdown
		log.Infow("shutdown the app is started", "signal", sig)
		defer log.Infow("shutdown the app is complete", "signal", sig)

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
