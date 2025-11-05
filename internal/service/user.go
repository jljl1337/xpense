package service

import (
	"context"

	"github.com/jljl1337/xpense/internal/crypto"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type UserService struct {
	queries *repository.Queries
}

func NewUserService(queries *repository.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

func (s *UserService) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	users, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return false, NewServiceErrorf(ErrCodeInternal, "failed to get user by username: %v", err)
	}

	if len(users) > 1 {
		return false, NewServiceError(ErrCodeInternal, "multiple users found with the same username")
	}

	return len(users) == 1, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*repository.User, error) {
	users, err := s.queries.GetUserByID(ctx, userID)
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

func (s *UserService) UpdateUsernameByID(ctx context.Context, userID, newUsername string) error {
	rows, err := s.queries.UpdateUserUsername(ctx, repository.UpdateUserUsernameParams{
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

func (s *UserService) UpdatePasswordByID(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Validate credentials
	users, err := s.queries.GetUserByID(ctx, userID)
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

	rows, err := s.queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
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

func (s *UserService) DeleteUserByID(ctx context.Context, userID string) error {
	rows, err := s.queries.DeleteUser(ctx, userID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete user: %v", err)
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no user deleted")
	}

	return nil
}
