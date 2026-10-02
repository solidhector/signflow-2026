// Package config loads runtime configuration from the environment.
//
// In development we lean on sensible defaults so you can `go run` the app with
// zero setup. In production every value is expected to come from the
// environment (Railway injects them). The struct GROWS as the app does — the
// database URL arrives in lesson 3, the session secret in lesson 4, and so on.

package config

import "os"

type Config struct {
	// Port is the TCP port the HTTP server listens on. Railway sets $PORT.
	Port string

	// Env is "dev" or "prod". Later it flips cookie Secure flags and the email
	// sender; for now it is just a label.
	Env string

	// BaseURL is the public origin used to build absolute links (in emails,
	// later). In dev it is http://localhost:PORT.
	BaseURL string
}

// IsProd reports whether we are running in production. One flag, read in many
// places later (Secure cookies, strict CSRF).
func (c Config) IsProd() bool { return c.Env == "prod" }

// Load reads configuration, applying dev-friendly defaults.
func Load() (Config, error) {
	cfg := Config{
		Port:    getenv("PORT", "8080"),
		Env:     getenv("APP_ENV", "dev"),
		BaseURL: os.Getenv("BASE_URL"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:" + cfg.Port
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
