package service

import (
	"context"
	"errors"
	"time"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type ExpenseService struct {
	queries *repository.Queries
}

func NewExpenseService(queries *repository.Queries) *ExpenseService {
	return &ExpenseService{
		queries: queries,
	}
}

// CreateExpense creates a new expense if the user has access to the book,
// category, and payment method.
//
// It returns true if the expense was created, false if the user has no access
// to the book, category, or payment method. It also returns false if the
// category or payment method does not belong to the book.
func (s *ExpenseService) CreateExpense(ctx context.Context, userID, bookID, categoryID, paymentMethodID string, date int64, amount float64, remark string) (bool, error) {
	// Check if the user has access to the book, category, and payment method
	canAccess, err := s.checkBookCategoryPaymentMethod(ctx, userID, bookID, categoryID, paymentMethodID)
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	// Create the expense
	currentTime := time.Now().UnixMilli()

	_, err = s.queries.CreateExpense(ctx, repository.CreateExpenseParams{
		ID:              generator.NewULID(),
		BookID:          bookID,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		Date:            date,
		Amount:          amount,
		Remark:          remark,
		CreatedAt:       currentTime,
		UpdatedAt:       currentTime,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetExpensesByBookID retrieves all expenses for a specific book with pagination.
//
// It returns an empty slice if no expenses are found in the book.
//
// It returns nil if the user has no access to the book or the book does not exist.
func (s *ExpenseService) GetExpensesByBookID(ctx context.Context, userID, bookID string, page int64, pageSize int64) ([]repository.Expense, error) {
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

	offset := (page - 1) * pageSize
	limit := pageSize
	expenses, err := s.queries.GetExpensesByBookID(ctx, repository.GetExpensesByBookIDParams{
		BookID: bookID,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetExpenseByID retrieves an expense by its ID if the user has access to the book.
//
// It returns nil if the expense does not exist or the user does not have access to the book.
func (s *ExpenseService) GetExpenseByID(ctx context.Context, userID, expenseID string) (*repository.Expense, error) {
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	if len(expenses) > 1 {
		return nil, errors.New("multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return nil, nil
	}

	expense := expenses[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: expense.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if !canAccess {
		return nil, nil
	}

	return &expense, nil
}

// UpdateExpense updates an existing expense if the user has access to the book,
// category, and payment method.
//
// It returns true if the expense was updated, false if the expense does not
// exist or the user has no access to the book, category, or payment method, or
// the category or payment method does not belong to the book.
func (s *ExpenseService) UpdateExpense(ctx context.Context, userID, expenseID, categoryID, paymentMethodID string, date int64, amount float64, remark string) (bool, error) {
	// Get the expense to find the book ID
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return false, err
	}

	if len(expenses) > 1 {
		return false, errors.New("multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return false, nil
	}

	expense := expenses[0]

	// Check if the user has access to the book, category, and payment method
	canAccess, err := s.checkBookCategoryPaymentMethod(ctx, userID, expense.BookID, categoryID, paymentMethodID)
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	// Update the expense
	currentTime := time.Now().UnixMilli()

	rows, err := s.queries.UpdateExpenseByID(ctx, repository.UpdateExpenseByIDParams{
		ID:              expenseID,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		Date:            date,
		Amount:          amount,
		Remark:          remark,
		UpdatedAt:       currentTime,
	})
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple expenses updated with the same ID")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}

// DeleteExpenseByID deletes an expense by its ID if the user has access to the
// book.
//
// It returns true if the expense was deleted, false if the expense does not
// exist or the user does not have access to the book.
func (s *ExpenseService) DeleteExpenseByID(ctx context.Context, userID, expenseID string) (bool, error) {
	// Get the expense to find the book ID
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return false, err
	}

	if len(expenses) > 1 {
		return false, errors.New("multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return false, nil
	}

	expense := expenses[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: expense.BookID,
		UserID: userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	// Proceed to delete the expense
	rows, err := s.queries.DeleteExpenseByID(ctx, expenseID)
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple expenses deleted with the same ID")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}

// checkBookCategoryPaymentMethod checks if the user has access to the book,
// category, and payment method.
//
// It also checks if the category and payment method belong to the book.
func (s *ExpenseService) checkBookCategoryPaymentMethod(ctx context.Context, userID, bookID, categoryID, paymentMethodID string) (bool, error) {
	// Check if the user has access to the book
	canAccessBook, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccessBook {
		return false, nil
	}

	// Check if the categories belongs to the book
	categories, err := s.queries.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return false, err
	}

	if len(categories) > 1 {
		return false, errors.New("multiple categories found with the same ID")
	}

	if len(categories) < 1 {
		return false, nil
	}

	category := categories[0]

	if category.BookID != bookID {
		return false, nil
	}

	// Check if the payment method belongs to the book
	paymentMethod, err := s.queries.GetPaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return false, err
	}

	if len(paymentMethod) > 1 {
		return false, errors.New("multiple payment methods found with the same ID")
	}

	if len(paymentMethod) < 1 {
		return false, nil
	}

	pm := paymentMethod[0]

	if pm.BookID != bookID {
		return false, nil
	}

	return true, nil
}
