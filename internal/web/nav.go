package web

// Nav carries the per-request bits the layout header needs: whether someone is
// signed in, and their email. A CSRF token joins it in lesson 6, once forms are
// protected.
type Nav struct {
	LoggedIn bool
	Email    string
}
