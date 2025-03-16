package middleware

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gorilla/sessions"
	"github.com/mavrk-mose/pay/config"
)

var store *sessions.CookieStore

// InitSessionStore initializes the session store with proper configuration
func InitSessionStore(cfg *config.Config) {
	// Generate a random 32-byte key for session encryption
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	// Create cookie store with the generated key
	store = sessions.NewCookieStore(key)

	// Configure session cookie
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.Session.Expire,
		HttpOnly: cfg.Cookie.HTTPOnly,
		Secure:   cfg.Cookie.Secure,
	}
}

// GetSessionStore returns the configured session store
func GetSessionStore() sessions.Store {
	return store
}

// GenerateSecureStateString generates a random state string for OAuth
func GenerateSecureStateString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
