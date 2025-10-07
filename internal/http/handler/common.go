package handler

import (
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
)

func NewActiveSessionCookie(sessionToken string) *http.Cookie {
	return &http.Cookie{
		Name:     env.SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Secure:   env.SessionCookieSecure,
		HttpOnly: env.SessionCookieHttpOnly,
		SameSite: env.SessionCookieSameSiteMode,
	}
}

func NewExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     env.SessionCookieName,
		Value:    "",
		Path:     "/",
		Secure:   env.SessionCookieSecure,
		HttpOnly: env.SessionCookieHttpOnly,
		SameSite: env.SessionCookieSameSiteMode,
		MaxAge:   -1,
	}
}
