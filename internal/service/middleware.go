package service

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/format"
	"github.com/jljl1337/xpense/internal/repository"
)

type MiddlewareService struct {
	db *sqlx.DB
}

func NewMiddlewareService(db *sqlx.DB) *MiddlewareService {
	return &MiddlewareService{
		db: db,
	}
}

// GetSessionUserIDAndRefreshSession validates the session token (and CSRF token),
// refreshes the session expiration, and returns the associated user ID.
func (s *MiddlewareService) GetSessionUserIDAndRefreshSession(ctx context.Context, sessionToken, CSRFToken string) (string, error) {
	queries := repository.New(s.db)

	sessions, err := queries.GetSessionByToken(ctx, sessionToken)

	if err != nil {
		return "", NewServiceErrorf(ErrCodeInternal, "failed to get session: %v", err)
	}

	if len(sessions) > 1 {
		return "", NewServiceError(ErrCodeInternal, "multiple sessions found with the same token")
	}

	if len(sessions) < 1 {
		return "", NewServiceError(ErrCodeUnauthorized, "unauthorized")
	}

	session := sessions[0]

	// Return unauthorized if the session is a pre session
	if !session.UserID.Valid {
		return "", NewServiceError(ErrCodeUnauthorized, "unauthorized")
	}

	// CSRF token does not match
	if CSRFToken != "" && session.CsrfToken != CSRFToken {
		return "", NewServiceError(ErrCodeUnauthorized, "unauthorized")
	}

	// Session expired
	now := time.Now()
	nowISO8601 := format.TimeToISO8601(now)
	if session.ExpiresAt < nowISO8601 {
		return "", NewServiceError(ErrCodeUnauthorized, "unauthorized")
	}

	// Only refresh session if remaining lifetime is below threshold
	expiresAt, err := format.ISO8601ToTime(session.ExpiresAt)
	if err != nil {
		return "", NewServiceErrorf(ErrCodeInternal, "failed to parse session expiration: %v", err)
	}

	remainingLifetimeMin := expiresAt.Sub(now).Minutes()
	if remainingLifetimeMin < float64(env.SessionRefreshThresholdMin) {
		newExpiresAt := format.TimeToISO8601(now.Add(time.Duration(env.SessionLifetimeMin) * time.Minute))
		rows, err := queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
			Token:     sessionToken,
			ExpiresAt: newExpiresAt,
			UpdatedAt: nowISO8601,
		})
		if err != nil {
			return "", NewServiceErrorf(ErrCodeInternal, "failed to refresh session: %v", err)
		}

		if rows < 1 {
			return "", NewServiceError(ErrCodeInternal, "no session updated")
		}
	}

	return session.UserID.String, nil
}
