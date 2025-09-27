package service

import (
	"context"
	"errors"
	"time"

	"github.com/jljl1337/xpense/internal/crypto"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/format"
	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type AuthService struct {
	queries *repository.Queries
}

func NewAuthService(queries *repository.Queries) *AuthService {
	return &AuthService{
		queries: queries,
	}
}

func (a *AuthService) SignUp(ctx context.Context, username, password string) error {
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	currentTime := generator.NowISO8601()

	_, err = a.queries.CreateUser(ctx, repository.CreateUserParams{
		ID:           generator.NewULID(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	})

	return err
}

// SignIn authenticates a user and creates a new session.
// It returns non-empty session token and CSRF token if the credentials are valid.
//
// If the credentials are invalid, it returns empty strings and no error.
//
// If an error occurs during the process, it returns the error.
func (a *AuthService) SignIn(ctx context.Context, username, password string) (string, string, error) {
	users, err := a.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if len(users) > 1 {
		return "", "", errors.New("multiple users found with the same username")
	}

	if len(users) < 1 {
		return "", "", nil
	}

	user := users[0]

	if !crypto.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", nil
	}

	sessionID := generator.NewULID()
	sessionToken := generator.NewToken(env.SessionTokenLength, env.SessionTokenCharset)
	CSRFToken := generator.NewToken(env.CSRFTokenLength, env.CSRFTokenCharset)
	currentTime := generator.NowISO8601()
	expiresAt := format.TimeToISO8601(time.Now().Add(24 * time.Hour))

	if _, err := a.queries.CreateSession(ctx, repository.CreateSessionParams{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     sessionToken,
		CsrfToken: CSRFToken,
		ExpiresAt: expiresAt,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}); err != nil {
		return "", "", err
	}

	return sessionToken, CSRFToken, nil
}

// GetSessionUserIDAndRefreshSession validates the session token and CSRF token,
// refreshes the session expiration, and returns the associated user ID.
//
// If the session is invalid or expired, it returns an empty string and no error.
//
// If an error occurs during the process, it returns the error.
func (a *AuthService) GetSessionUserIDAndRefreshSession(ctx context.Context, sessionToken, CSRFToken string) (string, error) {
	sessions, err := a.queries.GetSessionByToken(ctx, sessionToken)

	if err != nil {
		return "", err
	}

	if len(sessions) > 1 {
		return "", errors.New("multiple sessions found with the same token")
	}

	if len(sessions) < 1 {
		return "", nil
	}

	session := sessions[0]

	// CSRF token does not match
	if CSRFToken != "" && session.CsrfToken != CSRFToken {
		return "", nil
	}

	// Session expired
	now := time.Now()
	nowISO8601 := format.TimeToISO8601(now)
	if session.ExpiresAt < nowISO8601 {
		return "", nil
	}

	// Refresh the session expiration
	newExpiresAt := format.TimeToISO8601(now.Add(24 * time.Hour))
	rows, err := a.queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
		Token:     sessionToken,
		ExpiresAt: newExpiresAt,
		UpdatedAt: nowISO8601,
	})
	if err != nil {
		return "", err
	}

	if rows < 1 {
		return "", errors.New("no session updated")
	}

	return session.UserID, nil
}

func (a *AuthService) SignOut(ctx context.Context, sessionToken string) error {
	now := generator.NowISO8601()
	rows, err := a.queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
		Token:     sessionToken,
		ExpiresAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("no session updated")
	}

	return nil
}

func (a *AuthService) SignOutAllSession(ctx context.Context, userID string) error {
	now := generator.NowISO8601()
	rows, err := a.queries.UpdateSessionByUserID(ctx, repository.UpdateSessionByUserIDParams{
		UserID:    userID,
		ExpiresAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("no sessions deleted")
	}

	return nil
}

func (a *AuthService) CSRFToken(ctx context.Context, sessionToken string) (string, error) {
	sessions, err := a.queries.GetSessionByToken(ctx, sessionToken)

	if err != nil {
		return "", err
	}

	if len(sessions) > 1 {
		return "", errors.New("multiple sessions found with the same token")
	}

	if len(sessions) < 1 {
		return "", nil
	}

	session := sessions[0]

	return session.CsrfToken, nil
}
