// Package handlers holds the HTTP handlers. In this skeleton there is only the
// home page and a health check; later lessons add auth, documents, and signing.
package handlers

import (
	"log/slog"
	"net/http"

	"github.com/programuoki/signflow/internal/config"
	"github.com/programuoki/signflow/internal/db"
	"github.com/programuoki/signflow/internal/session"
	"github.com/programuoki/signflow/internal/web"
)

// Handlers bundles the dependencies every handler needs. Constructed once in
// main; its methods are registered on the router. It gains a Sessions manager
// here; the mailer arrives in lesson 7.
type Handlers struct {
	Cfg      config.Config
	Queries  *db.Queries
	Sessions *session.Manager
	Log      *slog.Logger
}

func New(cfg config.Config, q *db.Queries, sessions *session.Manager, log *slog.Logger) *Handlers {
	return &Handlers{Cfg: cfg, Queries: q, Sessions: sessions, Log: log}
}

// Home renders the landing page, including a live user count read from Postgres
// — proof the whole stack is wired end to end.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	count, err := h.Queries.CountUsers(r.Context())
	if err != nil {
		h.Log.Error("count users", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	render(w, r, http.StatusOK, web.Home(count))
}

// Health is a liveness probe for Railway / uptime checks. It also pings the DB,
// so an unreachable database reports unhealthy rather than silently serving 500s.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Queries.CountUsers(r.Context()); err != nil {
		http.Error(w, "db unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// serverError logs the underlying error and returns a generic 500 to the client.
func (h *Handlers) serverError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	h.Log.Error(msg, "err", err, "path", r.URL.Path)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
