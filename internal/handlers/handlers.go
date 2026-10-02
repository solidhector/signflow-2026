// Package handlers holds the HTTP handlers. In this skeleton there is only the
// home page and a health check; later lessons add auth, documents, and signing.
package handlers

import (
	"log/slog"
	"net/http"

	"github.com/programuoki/signflow/internal/config"
	"github.com/programuoki/signflow/internal/web"
)

// Handlers bundles the dependencies every handler needs. Constructed once in
// main; its methods are registered on the router. It gains fields as the app
// grows (a database in lesson 3, sessions in lesson 4, ...).
type Handlers struct {
	Cfg config.Config
	Log *slog.Logger
}

func New(cfg config.Config, log *slog.Logger) *Handlers {
	return &Handlers{Cfg: cfg, Log: log}
}

// Home renders the landing page from a type-safe Templ template.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, web.Home())
}

// Health is a plain-text liveness probe for Railway / uptime checks.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}
