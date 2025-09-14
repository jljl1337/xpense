package service

import (
	"context"
	"database/sql"
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
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*repository.User, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, userID string) error {
	rows, err := s.queries.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	if rows <= 0 {
		return errors.New("no user deleted")
	}

	return nil
}
