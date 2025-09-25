package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func (m *MiddlewareProvider) Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip for public routes
			if r.URL.Path == "/auth/sign-up" || r.URL.Path == "/auth/sign-in" || r.URL.Path == "/health" || r.URL.Path == "/users/exists" {
				next.ServeHTTP(w, r)
				return
			}

			// Get session token from cookie
			cookie, err := r.Cookie(env.SessionCookieName)
			if err != nil {
				// err is not nil only if the cookie is not present
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get CSRF token from header
			CSRFToken := r.Header.Get("X-CSRF-Token")

			if CSRFToken == "" && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch) {
				http.Error(w, "CSRF token is required", http.StatusUnauthorized)
				return
			}

			// Validate session token (and CSRF token)
			userID, err := m.authService.GetSessionUserIDAndRefreshSession(r.Context(), cookie.Value, CSRFToken)
			if err != nil {
				slog.Error("Failed to check session: " + err.Error())
				return
			}

			if userID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the context.
//
// It returns an error if the user ID is not found or is of an unexpected type.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("failed to get user ID from context")
	}
	return userID, nil
}
