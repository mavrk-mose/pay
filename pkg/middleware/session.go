package middleware

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gorilla/sessions"
	"github.com/mavrk-mose/pay/config"
)

var store *sessions.CookieStore

func InitSessionStore(cfg *config.Config) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	store = sessions.NewCookieStore(key)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.Session.Expire,
		HttpOnly: cfg.Cookie.HTTPOnly,
		Secure:   cfg.Cookie.Secure,
	}
}

func GetSessionStore() sessions.Store {
	return store
}

func GenerateSecureStateString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
