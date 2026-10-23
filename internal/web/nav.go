package web

// Nav carries the per-request bits the layout header needs: whether someone is
// signed in, their email, and the CSRF token (for the logout form, which is a
// POST, and for any other form the page renders).
type Nav struct {
	LoggedIn  bool
	Email     string
	CSRFToken string
}
