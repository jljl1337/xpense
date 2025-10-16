package service

import (
	"context"
	"errors"

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
		return false, err
	}

	if len(users) > 1 {
		return false, errors.New("multiple users found with the same username")
	}

	return len(users) == 1, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*repository.User, error) {
	users, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(users) > 1 {
		return nil, errors.New("multiple users found with the same ID")
	}

	if len(users) < 1 {
		return nil, nil
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
		return err
	}

	if rows < 1 {
		return errors.New("no user updated")
	}
	if rows > 1 {
		return errors.New("multiple users updated")
	}

	return nil
}

func (s *UserService) UpdatePasswordByID(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Validate credentials
	users, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if len(users) > 1 {
		return errors.New("multiple users found with the same ID")
	}

	if len(users) < 1 {
		return errors.New("user not found")
	}

	user := users[0]

	if !crypto.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("old password is incorrect")
	}

	// Update password hash
	passwordHash, err := crypto.HashPassword(newPassword, env.PasswordBcryptCost)
	if err != nil {
		return err
	}

	rows, err := s.queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		PasswordHash: passwordHash,
		UpdatedAt:    generator.NowISO8601(),
		ID:           userID,
	})
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("no user updated")
	}
	if rows > 1 {
		return errors.New("multiple users updated")
	}

	return nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, userID string) error {
	rows, err := s.queries.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("no user deleted")
	}

	return nil
}
