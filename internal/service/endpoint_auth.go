package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/jljl1337/xpense/internal/crypto"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/format"
	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

func (s *EndpointService) SignUp(ctx context.Context, username, password string) error {
	queries := repository.New(s.db)

	users, err := queries.GetUserByUsername(ctx, username)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get user by username: %v", err)
	}

	if len(users) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple users found with the same username")
	}

	if len(users) > 0 {
		return NewServiceError(ErrCodeConflict, "username already exists")
	}

	passwordHash, err := crypto.HashPassword(password, env.PasswordBcryptCost)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to hash password: %v", err)
	}

	currentTime := generator.NowISO8601()

	if _, err = queries.CreateUser(ctx, repository.CreateUserParams{
		ID:           generator.NewULID(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}); err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to create user: %v", err)
	}

	return nil
}

// GetPreSession creates a pre-session with no associated user.
// It returns a non-empty session token and CSRF token.
func (s *EndpointService) GetPreSession(ctx context.Context) (string, string, error) {
	queries := repository.New(s.db)

	sessionID := generator.NewULID()
	sessionToken := generator.NewToken(env.SessionTokenLength, env.SessionTokenCharset)
	CSRFToken := generator.NewToken(env.CSRFTokenLength, env.CSRFTokenCharset)
	currentTime := generator.NowISO8601()
	expiresAt := format.TimeToISO8601(time.Now().Add(time.Duration(env.PreSessionLifetimeMin) * time.Minute))

	if _, err := queries.CreateSession(ctx, repository.CreateSessionParams{
		ID:        sessionID,
		UserID:    sql.NullString{Valid: false},
		Token:     sessionToken,
		CsrfToken: CSRFToken,
		ExpiresAt: expiresAt,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}); err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to create pre-session: %v", err)
	}

	return sessionToken, CSRFToken, nil
}

// SignIn authenticates a user and creates a new session.
// It returns non-empty session token and CSRF token if the credentials are valid.
func (s *EndpointService) SignIn(ctx context.Context, preSessionToken, preSessionCSRFToken, username, password string) (string, string, error) {
	queries := repository.New(s.db)

	// Validate pre-session
	sessions, err := queries.GetSessionByToken(ctx, preSessionToken)

	if err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to get pre-session: %v", err)
	}

	if len(sessions) > 1 {
		return "", "", NewServiceError(ErrCodeInternal, "multiple sessions found with the same token")
	}

	if len(sessions) < 1 {
		return "", "", NewServiceError(ErrCodeUnauthorized, "invalid credentials")
	}

	session := sessions[0]

	// Check if the session is already associated with a user
	if session.UserID.Valid {
		return "", "", NewServiceError(ErrCodeUnauthorized, "invalid credentials")
	}

	// CSRF token does not match
	if preSessionCSRFToken != "" && session.CsrfToken != preSessionCSRFToken {
		return "", "", NewServiceError(ErrCodeUnauthorized, "invalid credentials")
	}

	// Session expired
	now := time.Now()
	nowISO8601 := format.TimeToISO8601(now)
	if session.ExpiresAt < nowISO8601 {
		return "", "", NewServiceError(ErrCodeUnauthorized, "session expired")
	}

	// Validate credentials
	users, err := queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to get user by username: %v", err)
	}

	if len(users) > 1 {
		return "", "", NewServiceError(ErrCodeInternal, "multiple users found with the same username")
	}

	if len(users) < 1 {
		return "", "", NewServiceError(ErrCodeUnauthorized, "invalid credentials")
	}

	user := users[0]

	if !crypto.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", NewServiceError(ErrCodeUnauthorized, "invalid credentials")
	}

	cost, err := crypto.Cost(user.PasswordHash)
	if err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to get password hash cost: %v", err)
	}

	// Rehash password if the cost is lower than the current standard
	currentTime := generator.NowISO8601()

	if cost < env.PasswordBcryptCost {
		newHash, err := crypto.HashPassword(password, env.PasswordBcryptCost)
		if err != nil {
			return "", "", NewServiceErrorf(ErrCodeInternal, "failed to hash password: %v", err)
		}

		rows, err := queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
			PasswordHash: newHash,
			UpdatedAt:    currentTime,
			ID:           user.ID,
		})
		if err != nil {
			return "", "", NewServiceErrorf(ErrCodeInternal, "failed to update user password hash: %v", err)
		} else if rows < 1 {
			return "", "", NewServiceError(ErrCodeInternal, "no user updated with the new password hash")
		} else if rows > 1 {
			return "", "", NewServiceError(ErrCodeInternal, "multiple users updated with the same ID")
		}
	}

	sessionID := generator.NewULID()
	sessionToken := generator.NewToken(env.SessionTokenLength, env.SessionTokenCharset)
	CSRFToken := generator.NewToken(env.CSRFTokenLength, env.CSRFTokenCharset)
	expiresAt := format.TimeToISO8601(time.Now().Add(time.Duration(env.SessionLifetimeMin) * time.Hour))

	// Deactivate the pre-session
	rows, err := queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
		Token:     preSessionToken,
		ExpiresAt: nowISO8601,
		UpdatedAt: nowISO8601,
	})
	if err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to update pre-session: %v", err)
	} else if rows < 1 {
		return "", "", NewServiceError(ErrCodeInternal, "no pre-session updated")
	} else if rows > 1 {
		return "", "", NewServiceError(ErrCodeInternal, "multiple pre-sessions updated with the same token")
	}

	// Create a new session associated with the user
	rows, err = queries.CreateSession(ctx, repository.CreateSessionParams{
		ID:        sessionID,
		UserID:    sql.NullString{String: user.ID, Valid: true},
		Token:     sessionToken,
		CsrfToken: CSRFToken,
		ExpiresAt: expiresAt,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	})
	if err != nil {
		return "", "", NewServiceErrorf(ErrCodeInternal, "failed to create session: %v", err)
	} else if rows < 1 {
		return "", "", NewServiceError(ErrCodeInternal, "no session created")
	} else if rows > 1 {
		return "", "", NewServiceError(ErrCodeInternal, "multiple sessions created with the same ID")
	}

	return sessionToken, CSRFToken, nil
}

func (s *EndpointService) SignOut(ctx context.Context, sessionToken string) error {
	queries := repository.New(s.db)

	now := generator.NowISO8601()
	rows, err := queries.UpdateSessionByToken(ctx, repository.UpdateSessionByTokenParams{
		Token:     sessionToken,
		ExpiresAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to sign out session: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no session updated")
	}

	return nil
}

func (s *EndpointService) SignOutAllSession(ctx context.Context, userID string) error {
	queries := repository.New(s.db)

	now := generator.NowISO8601()
	rows, err := queries.UpdateSessionByUserID(ctx, repository.UpdateSessionByUserIDParams{
		UserID:    sql.NullString{String: userID, Valid: true},
		ExpiresAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to sign out all sessions: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no sessions deleted")
	}

	return nil
}

func (s *EndpointService) CSRFToken(ctx context.Context, sessionToken string) (string, error) {
	queries := repository.New(s.db)

	sessions, err := queries.GetSessionByToken(ctx, sessionToken)

	if err != nil {
		return "", NewServiceErrorf(ErrCodeInternal, "failed to get session: %v", err)
	}

	if len(sessions) > 1 {
		return "", NewServiceError(ErrCodeInternal, "multiple sessions found with the same token")
	}

	if len(sessions) < 1 {
		return "", NewServiceError(ErrCodeUnauthorized, "invalid session")
	}

	session := sessions[0]

	return session.CsrfToken, nil
}
