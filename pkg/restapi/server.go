package restapi

import (
	"log"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
)

func CreateAndRunServer(mux *chi.Mux, serverAddr string) error {
	s := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}

	log.Printf("Starting server at PORT %s\n", serverAddr)

	return s.ListenAndServe()
}
