package service

import (
	"context"
	"errors"
	"time"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type CategoryService struct {
	queries *repository.Queries
}

func NewCategoryService(queries *repository.Queries) *CategoryService {
	return &CategoryService{
		queries: queries,
	}
}

// CreateCategory creates a new category if the user has access to the book.
//
// It returns true if the category was created, false if the user has no access to the book
func (s *CategoryService) CreateCategory(ctx context.Context, userID, bookID, name, description string) (bool, error) {
	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	currentTime := time.Now().UnixMilli()

	_, err = s.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:          generator.NewULID(),
		BookID:      bookID,
		Name:        name,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetCategoriesByBookID retrieves all categories for a specific book.
//
// It returns an empty slice if no categories are found in the book.
//
// It returns nil if the user has no access to the book.
func (s *CategoryService) GetCategoriesByBookID(ctx context.Context, userID, bookID string) ([]repository.Category, error) {
	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if !canAccess {
		return nil, nil
	}

	categories, err := s.queries.GetCategoriesByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID if the user has access to the book.
//
// It returns nil if the category does not exist or the user does not have
// access to the book.
func (s *CategoryService) GetCategoryByID(ctx context.Context, userID, categoryID string) (*repository.Category, error) {
	// Get the category to find the book ID
	categories, err := s.queries.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if len(categories) > 1 {
		return nil, errors.New("multiple categories found with the same ID")
	}

	if len(categories) < 1 {
		return nil, nil
	}

	category := categories[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: category.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if !canAccess {
		return nil, nil
	}

	return &category, nil
}

// UpdateCategoryByID updates a category if the user has access to the book.
//
// It returns true if the category was updated, false if the user has no access
// to the category or the category does not exist
func (s *CategoryService) UpdateCategoryByID(ctx context.Context, userID, categoryID, name, description string) (bool, error) {
	// Check if the user has access to the category
	canAccess, err := s.queries.CheckCategoryAccess(ctx, repository.CheckCategoryAccessParams{
		CategoryID: categoryID,
		UserID:     userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	currentTime := time.Now().UnixMilli()

	rows, err := s.queries.UpdateCategoryByID(ctx, repository.UpdateCategoryByIDParams{
		ID:          categoryID,
		Name:        name,
		Description: description,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple categories updated, data integrity issue")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}

// DeleteCategoryByID deletes a category if the user has access to the book.
//
// It returns true if the category was deleted, false if the user has no access
// to the category or the category does not exist
func (s *CategoryService) DeleteCategoryByID(ctx context.Context, userID, categoryID string) (bool, error) {
	// Check if the user has access to the category
	canAccess, err := s.queries.CheckCategoryAccess(ctx, repository.CheckCategoryAccessParams{
		CategoryID: categoryID,
		UserID:     userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	rows, err := s.queries.DeleteCategoryByID(ctx, categoryID)
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple categories deleted, data integrity issue")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}
