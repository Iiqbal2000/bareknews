package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App
type App struct {
	mux      *chi.Mux
	shutdown chan os.Signal
	middlewares []Middleware
}

// NewApp 
func NewApp(shutdown chan os.Signal, m ...Middleware) App {
	mux := chi.NewMux()

	return App{
		mux:      mux,
		shutdown: shutdown,
		middlewares: m,
	}
}

// ServeHTTP implements Handler interface.
func (app App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(rw, r)
}

// Handle 
func (app App) Handle(method, pattern string, handler Handler) {

	handler = SetMiddlewares(app.middlewares, handler)
	
	// The function executed for each request.
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		ctx := r.Context()

		// Capture the parent request span from the context.
		span := trace.SpanFromContext(ctx)

		traceID := span.SpanContext().TraceID().String()
		timenow := time.Now().UTC().String()

		span.SetAttributes(
			attribute.String("TraceID", traceID),
			attribute.String("time", timenow),
		)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			app.signalShutdown()
			return
		}
	})
	
	app.mux.Method(method, pattern, otelhttp.WithRouteTag(pattern, h))
}

// signalShutdown emits a signal to shutdown the app.
func (app App) signalShutdown() {
	app.shutdown <- syscall.SIGTERM
}
