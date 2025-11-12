package service

import (
	"context"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

func (s *EndpointService) CreateBook(ctx context.Context, userID, name, description string) error {
	queries := repository.New(s.db)

	currentTime := generator.NowISO8601()

	_, err := queries.CreateBook(ctx, repository.CreateBookParams{
		ID:          generator.NewULID(),
		UserID:      userID,
		Name:        name,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to create book: %v", err)
	}

	return nil
}

func (s *EndpointService) GetBooksCountByUserID(ctx context.Context, userID string) (int64, error) {
	queries := repository.New(s.db)

	countResult, err := queries.GetBooksCountByUserID(ctx, userID)
	if err != nil {
		return 0, NewServiceErrorf(ErrCodeInternal, "failed to get books count: %v", err)
	}

	return countResult, nil
}

// GetBooksByUserID retrieves a paginated list of books for a specific user.
//
// It returns an empty slice if no books are found.
func (s *EndpointService) GetBooksByUserID(ctx context.Context, userID string, page int64, pageSize int64) ([]repository.Book, error) {
	queries := repository.New(s.db)

	offset := (page - 1) * pageSize
	limit := pageSize
	books, err := queries.GetBooksByUserID(ctx, repository.GetBooksByUserIDParams{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get books by user ID: %v", err)
	}

	return books, nil
}

// GetBookByID retrieves a book by its ID if the user has access to it.
func (s *EndpointService) GetBookByID(ctx context.Context, userID, bookID string) (*repository.Book, error) {
	queries := repository.New(s.db)

	// Check if the user has access to the book
	canAccess, err := queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return nil, NewServiceErrorf(ErrCodeNotFound, "book not found or access denied")
	}

	// Fetch the book details
	books, err := queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get book by ID: %v", err)
	}

	if len(books) > 1 {
		return nil, NewServiceError(ErrCodeInternal, "multiple books found with the same ID")
	}

	if len(books) < 1 {
		return nil, NewServiceErrorf(ErrCodeNotFound, "book not found or access denied")
	}

	book := books[0]

	return &book, nil
}

// UpdateBookByID updates a book's name and description if the user has access to it.
func (s *EndpointService) UpdateBookByID(ctx context.Context, userID, bookID, name, description string) error {
	queries := repository.New(s.db)

	// Check if the user has access to the book
	canAccess, err := queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "book not found or access denied")
	}

	// Proceed to update the book
	rows, err := queries.UpdateBookByID(ctx, repository.UpdateBookByIDParams{
		ID:          bookID,
		Name:        name,
		Description: description,
		UpdatedAt:   generator.NowISO8601(),
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update book: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple books updated, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no book updated")
	}

	return nil
}

// DeleteBookByID deletes a book by its ID if the user has access to it.
func (s *EndpointService) DeleteBookByID(ctx context.Context, userID, bookID string) error {
	queries := repository.New(s.db)

	// Check if the user has access to the book
	canAccess, err := queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "book not found or access denied")
	}

	// Proceed to delete the book
	rows, err := queries.DeleteBookByID(ctx, bookID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete book: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple books deleted, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no book deleted")
	}

	return nil
}
