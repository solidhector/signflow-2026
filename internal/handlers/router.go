package handlers

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router builds the Chi router with the middleware stack and all routes.
func (h *Handlers) Router(staticFS fs.FS) http.Handler {
	r := chi.NewRouter()

	// Base middleware, outermost first.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)    // request logging
	r.Use(middleware.Recoverer) // turn panics into 500s instead of crashing

	// Resolve the session cookie into the request context for every route.
	r.Use(h.LoadUser)

	// Static assets (CSS, vendored htmx), served from the embedded FS.
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Public routes.
	r.Get("/", h.Home)
	r.Get("/healthz", h.Health)

	// Auth routes.
	r.Get("/register", h.RegisterForm)
	r.Post("/register", h.Register)
	r.Get("/login", h.LoginForm)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)

	// Authenticated routes.
	r.Group(func(pr chi.Router) {
		pr.Use(h.RequireAuth)
		pr.Get("/dashboard", h.Dashboard)
		pr.Get("/account", h.Account)
	})

	return r
}
