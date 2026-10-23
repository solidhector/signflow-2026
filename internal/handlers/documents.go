package handlers

import (
	"net/http"
	"strings"

	"github.com/programuoki/signflow/internal/db"
	"github.com/programuoki/signflow/internal/web"
)

// Dashboard lists the signed-in user's documents.
func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := mustUser(r)
	docs, err := h.Queries.ListDocumentsByOwner(r.Context(), user.ID)
	if err != nil {
		h.serverError(w, r, "list documents", err)
		return
	}
	render(w, r, http.StatusOK, web.Dashboard(h.nav(r), viewDocuments(docs)))
}

// UploadDocument accepts a multipart file, hashes and stores it, and records it
// as a draft owned by the current user.
func (h *Handlers) UploadDocument(w http.ResponseWriter, r *http.Request) {
	user := mustUser(r)

	// Cap the request body. +1KiB slack covers multipart field overhead so a
	// file exactly at the limit still succeeds.
	r.Body = http.MaxBytesReader(w, r.Body, h.Cfg.MaxUploadBytes+1024)

	file, header, err := r.FormFile("file")
	if err != nil {
		// Covers both "no file" and "too large" (MaxBytesReader error).
		h.Log.Warn("upload rejected", "err", err, "user", user.Email)
		http.Error(w, "upload failed: file missing or too large", http.StatusRequestEntityTooLarge)
		return
	}
	defer file.Close()

	if header.Size == 0 {
		http.Error(w, "upload failed: empty file", http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// The store hashes while it writes — the file is never fully in memory.
	key, size, hash, err := h.Store.Save(r.Context(), file)
	if err != nil {
		h.serverError(w, r, "store file", err)
		return
	}

	doc, err := h.Queries.CreateDocument(r.Context(), db.CreateDocumentParams{
		OwnerID:     user.ID,
		Filename:    sanitizeFilename(header.Filename),
		ContentType: contentType,
		Size:        size,
		FileHash:    hash,
		StorageKey:  key,
	})
	if err != nil {
		_ = h.Store.Delete(r.Context(), key) // don't orphan the blob
		h.serverError(w, r, "create document", err)
		return
	}
	h.Log.Info("document uploaded", "id", doc.ID.String(), "owner", user.Email, "size", size, "sha256", hash)

	h.respondWithDocList(w, r, user)
}

// respondWithDocList returns the refreshed list partial to HTMX callers, or
// redirects to the dashboard for a plain form post.
func (h *Handlers) respondWithDocList(w http.ResponseWriter, r *http.Request, user db.User) {
	if r.Header.Get("HX-Request") == "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	docs, err := h.Queries.ListDocumentsByOwner(r.Context(), user.ID)
	if err != nil {
		h.serverError(w, r, "list documents", err)
		return
	}
	render(w, r, http.StatusOK, web.DocumentList(viewDocuments(docs)))
}

func mustUser(r *http.Request) db.User {
	u, _ := UserFrom(r.Context()) // routes are behind RequireAuth
	return u
}

// sanitizeFilename strips path separators and control characters so a filename
// can't traverse directories or inject headers.
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.Map(func(rn rune) rune {
		if rn < 0x20 || rn == 0x7f || rn ==
			'"' {
			return -1
		}
		return rn
	}, name)
	name = strings.TrimSpace(name)
	if name == "" {
		return "document"
	}
	return name
}
