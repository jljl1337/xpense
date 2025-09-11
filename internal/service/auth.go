package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jljl1337/xpense/internal/crypto"
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

func (a *AuthService) SignUp(email, password string) error {
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	ctx := context.Background()
	currentTime := time.Now().UnixMilli()

	_, err = a.queries.CreateUser(ctx, repository.CreateUserParams{
		ID:           generator.NewKSUID(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	})

	return err
}

// Login authenticates a user and creates a new session.
// It returns non-empty session token and CSRF token if the credentials are valid.
//
// If the credentials are invalid, it returns empty strings and no error.
//
// If an error occurs during the process, it returns the error.
func (a *AuthService) Login(email, password string) (string, string, error) {
	ctx := context.Background()

	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		return "", "", err
	}

	if !crypto.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", nil
	}

	sessionID := generator.NewKSUID()
	sessionToken := generator.NewToken(16)
	CSRFToken := generator.NewToken(16)
	currentTime := time.Now().UnixMilli()
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

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
func (a *AuthService) GetSessionUserIDAndRefreshSession(sessionToken, CSRFToken string) (string, error) {
	ctx := context.Background()

	session, err := a.queries.GetSessionByToken(ctx, sessionToken)

	if err != nil {
		// No session with the given token
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	// CSRF token does not match
	if CSRFToken != "" && session.CsrfToken != CSRFToken {
		return "", nil
	}

	// Session expired
	now := time.Now()
	nowMillis := now.UnixMilli()
	if session.ExpiresAt < nowMillis {
		return "", nil
	}

	// Refresh the session expiration
	newExpiresAt := now.Add(24 * time.Hour).UnixMilli()
	rows, err := a.queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
		Token:     sessionToken,
		ExpiresAt: newExpiresAt,
		UpdatedAt: nowMillis,
	})
	if err != nil {
		return "", err
	}

	if rows < 1 {
		return "", errors.New("no session updated")
	}

	return session.UserID, nil
}

func (a *AuthService) Logout(sessionToken string) error {
	ctx := context.Background()
	now := time.Now().UnixMilli()
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
