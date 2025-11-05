package service

import (
	"context"

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
func (s *CategoryService) CreateCategory(ctx context.Context, userID, bookID, name, description string) error {
	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeUnprocessable, "book not found or access denied")
	}

	currentTime := generator.NowISO8601()

	_, err = s.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:          generator.NewULID(),
		BookID:      bookID,
		Name:        name,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to create category: %v", err)
	}

	return nil
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
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return nil, NewServiceError(ErrCodeUnprocessable, "book not found or access denied")
	}

	categories, err := s.queries.GetCategoriesByBookID(ctx, bookID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get categories by book ID: %v", err)
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
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get category by ID: %v", err)
	}

	if len(categories) > 1 {
		return nil, NewServiceError(ErrCodeInternal, "multiple categories found with the same ID")
	}

	if len(categories) < 1 {
		return nil, NewServiceError(ErrCodeNotFound, "category not found or access denied")
	}

	category := categories[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: category.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return nil, NewServiceError(ErrCodeNotFound, "book not found or access denied")
	}

	return &category, nil
}

// UpdateCategoryByID updates a category if the user has access to the book.
//
// It returns true if the category was updated, false if the user has no access
// to the category or the category does not exist
func (s *CategoryService) UpdateCategoryByID(ctx context.Context, userID, categoryID, name, description string) error {
	// Check if the user has access to the category
	canAccess, err := s.queries.CheckCategoryAccess(ctx, repository.CheckCategoryAccessParams{
		CategoryID: categoryID,
		UserID:     userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check category access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "category not found or access denied")
	}

	rows, err := s.queries.UpdateCategoryByID(ctx, repository.UpdateCategoryByIDParams{
		ID:          categoryID,
		Name:        name,
		Description: description,
		UpdatedAt:   generator.NowISO8601(),
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update category: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple categories updated, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no category updated")
	}

	return nil
}

// DeleteCategoryByID deletes a category if the user has access to the book.
//
// It returns true if the category was deleted, false if the user has no access
// to the category or the category does not exist
func (s *CategoryService) DeleteCategoryByID(ctx context.Context, userID, categoryID string) error {
	// Check if the user has access to the category
	canAccess, err := s.queries.CheckCategoryAccess(ctx, repository.CheckCategoryAccessParams{
		CategoryID: categoryID,
		UserID:     userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check category access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "category not found or access denied")
	}

	rows, err := s.queries.DeleteCategoryByID(ctx, categoryID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete category: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple categories deleted, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no category deleted")
	}

	return nil
}
