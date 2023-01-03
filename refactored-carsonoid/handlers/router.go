package handlers

import (
	"net/http"
	"time"

	"github.com/bradleyshawkins/go-clean-architecture/refactored-carsonoid/integrations"
	"github.com/go-chi/chi/v5"
)

// NewMux creates a new chi.Mux and registers the handlers
func NewMux() *chi.Mux {
	m := chi.NewMux()

	web := integrations.NewWebIntegration(&http.Client{
		Timeout: 10 * time.Second,
	})
	noop := integrations.NewNoopIntegration()

	m.Post("/download", NewFileDownloader(web, noop))

	return m
}
