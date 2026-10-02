package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

// render writes a templ component as the HTTP response with the given status.
// Centralising this keeps handlers terse and gives us one place to set headers.
func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	// A render error mid-stream can't change the status code (already written);
	// there's nothing useful to do but drop the connection, which Render does.
	_ = c.Render(r.Context(), w)
}
