package service

import (
	"context"
	"errors"

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

func (s *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	users, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if len(users) > 1 {
		return false, errors.New("multiple users found with the same email")
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
