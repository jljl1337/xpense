package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type BookService struct {
	queries *repository.Queries
}

func NewBookService(queries *repository.Queries) *BookService {
	return &BookService{
		queries: queries,
	}
}

func (s *BookService) CreateBook(ctx context.Context, userID, name, description string) error {
	currentTime := generator.NowISO8601()

	_, err := s.queries.CreateBook(ctx, repository.CreateBookParams{
		ID:          generator.NewULID(),
		UserID:      userID,
		Name:        name,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetBooksByUserID retrieves a paginated list of books for a specific user.
//
// It returns an empty slice if no books are found.
func (s *BookService) GetBooksByUserID(ctx context.Context, userID string, page int64, pageSize int64) ([]repository.Book, error) {
	offset := (page - 1) * pageSize
	limit := pageSize
	books, err := s.queries.GetBooksByUserID(ctx, repository.GetBooksByUserIDParams{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetBookByID retrieves a book by its ID if the user has access to it.
//
// It returns nil if the book does not exist or the user does not have access to it.
func (s *BookService) GetBookByID(ctx context.Context, userID, bookID string) (*repository.Book, error) {
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

	// Fetch the book details
	books, err := s.queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if len(books) > 1 {
		return nil, errors.New("multiple books found with the same ID")
	}

	if len(books) < 1 {
		return nil, sql.ErrNoRows
	}

	book := books[0]

	return &book, nil
}

// UpdateBookByID updates a book's name and description if the user has access to it.
//
// It returns true if the update was successful, false if the book does not exist
// or the user does not have access to it.
func (s *BookService) UpdateBookByID(ctx context.Context, userID, bookID, name, description string) (bool, error) {
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

	// Proceed to update the book
	rows, err := s.queries.UpdateBookByID(ctx, repository.UpdateBookByIDParams{
		ID:          bookID,
		Name:        name,
		Description: description,
		UpdatedAt:   generator.NowISO8601(),
	})
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple books updated, data integrity issue")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}

// DeleteBookByID deletes a book by its ID if the user has access to it.
//
// It returns true if the deletion was successful, false if the book does not exist
// or the user does not have access to it.
func (s *BookService) DeleteBookByID(ctx context.Context, userID, bookID string) (bool, error) {
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

	// Proceed to delete the book
	rows, err := s.queries.DeleteBookByID(ctx, bookID)
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple books deleted, data integrity issue")
	}

	if rows < 1 {
		return false, errors.New("no book deleted")
	}

	return true, nil
}
