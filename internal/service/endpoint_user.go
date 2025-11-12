package service

import (
	"context"

	"github.com/jljl1337/xpense/internal/crypto"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

func (s *EndpointService) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	queries := repository.New(s.db)

	users, err := queries.GetUserByUsername(ctx, username)
	if err != nil {
		return false, NewServiceErrorf(ErrCodeInternal, "failed to get user by username: %v", err)
	}

	if len(users) > 1 {
		return false, NewServiceError(ErrCodeInternal, "multiple users found with the same username")
	}

	return len(users) == 1, nil
}

func (s *EndpointService) GetUserByID(ctx context.Context, userID string) (*repository.User, error) {
	queries := repository.New(s.db)

	users, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get user: %v", err)
	}

	if len(users) > 1 {
		return nil, NewServiceError(ErrCodeInternal, "multiple users found with the same ID")
	}

	if len(users) < 1 {
		return nil, NewServiceError(ErrCodeNotFound, "user not found")
	}

	return &users[0], nil
}

func (s *EndpointService) UpdateUsernameByID(ctx context.Context, userID, newUsername string) error {
	queries := repository.New(s.db)

	rows, err := queries.UpdateUserUsername(ctx, repository.UpdateUserUsernameParams{
		ID:        userID,
		Username:  newUsername,
		UpdatedAt: generator.NowISO8601(),
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update username: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no user updated")
	}
	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple users updated")
	}

	return nil
}

func (s *EndpointService) UpdatePasswordByID(ctx context.Context, userID, oldPassword, newPassword string) error {
	queries := repository.New(s.db)

	// Validate credentials
	users, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get user: %v", err)
	}

	if len(users) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple users found with the same ID")
	}

	if len(users) < 1 {
		return NewServiceError(ErrCodeNotFound, "user not found")
	}

	user := users[0]

	if !crypto.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return NewServiceError(ErrCodeUnprocessable, "old password is incorrect")
	}

	// Update password hash
	passwordHash, err := crypto.HashPassword(newPassword, env.PasswordBcryptCost)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to hash password: %v", err)
	}

	rows, err := queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		PasswordHash: passwordHash,
		UpdatedAt:    generator.NowISO8601(),
		ID:           userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update password: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no user updated")
	}
	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple users updated")
	}

	return nil
}

func (s *EndpointService) DeleteUserByID(ctx context.Context, userID string) error {
	queries := repository.New(s.db)

	rows, err := queries.DeleteUser(ctx, userID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete user: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no user deleted")
	}

	return nil
}
