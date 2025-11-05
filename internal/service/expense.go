package service

import (
	"context"

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
func (s *ExpenseService) CreateExpense(ctx context.Context, userID, bookID, categoryID, paymentMethodID, date string, amount float64, remark string) error {
	// Check if the user has access to the book, category, and payment method
	err := s.checkBookCategoryPaymentMethod(ctx, userID, bookID, categoryID, paymentMethodID)
	if err != nil {
		return err
	}

	// Create the expense
	currentTime := generator.NowISO8601()

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
		return NewServiceErrorf(ErrCodeInternal, "failed to create expense: %v", err)
	}

	return nil
}

func (s *ExpenseService) GetExpensesCountByBookID(ctx context.Context, userID, bookID, categoryID, paymentMethodID, remark string) (int64, error) {
	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return 0, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return 0, NewServiceError(ErrCodeNotFound, "book not found or access denied")
	}

	countResult, err := s.queries.GetExpenseCountByBookID(ctx, repository.GetExpenseCountByBookIDParams{
		BookID:          bookID,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		Remark:          remark,
	})
	if err != nil {
		return 0, NewServiceErrorf(ErrCodeInternal, "failed to get expenses count: %v", err)
	}

	return countResult, nil
}

// GetExpensesByBookID retrieves all expenses for a specific book with pagination.
//
// It returns an empty slice if no expenses are found in the book.
//
// It returns nil if the user has no access to the book or the book does not exist.
func (s *ExpenseService) GetExpensesByBookID(ctx context.Context, userID, bookID, categoryID, paymentMethodID, remark string, page int64, pageSize int64) ([]repository.Expense, error) {
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

	offset := (page - 1) * pageSize
	limit := pageSize
	expenses, err := s.queries.GetExpensesByBookID(ctx, repository.GetExpensesByBookIDParams{
		BookID:          bookID,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		Remark:          remark,
		Offset:          offset,
		Limit:           limit,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get expenses by book ID: %v", err)
	}

	return expenses, nil
}

// GetExpenseByID retrieves an expense by its ID if the user has access to the book.
//
// It returns nil if the expense does not exist or the user does not have access to the book.
func (s *ExpenseService) GetExpenseByID(ctx context.Context, userID, expenseID string) (*repository.Expense, error) {
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get expense by ID: %v", err)
	}

	if len(expenses) > 1 {
		return nil, NewServiceError(ErrCodeInternal, "multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return nil, NewServiceError(ErrCodeNotFound, "expense not found or access denied")
	}

	expense := expenses[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: expense.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return nil, NewServiceError(ErrCodeNotFound, "expense not found or access denied")
	}

	return &expense, nil
}

// UpdateExpense updates an existing expense if the user has access to the book,
// category, and payment method.
//
// It returns true if the expense was updated, false if the expense does not
// exist or the user has no access to the book, category, or payment method, or
// the category or payment method does not belong to the book.
func (s *ExpenseService) UpdateExpense(ctx context.Context, userID, expenseID, categoryID, paymentMethodID, date string, amount float64, remark string) error {
	// Get the expense to find the book ID
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get expense: %v", err)
	}

	if len(expenses) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return NewServiceError(ErrCodeNotFound, "expense not found or access denied")
	}

	expense := expenses[0]

	// Check if the user has access to the book, category, and payment method
	err = s.checkBookCategoryPaymentMethod(ctx, userID, expense.BookID, categoryID, paymentMethodID)
	if err != nil {
		return err
	}

	// Update the expense
	rows, err := s.queries.UpdateExpenseByID(ctx, repository.UpdateExpenseByIDParams{
		ID:              expenseID,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		Date:            date,
		Amount:          amount,
		Remark:          remark,
		UpdatedAt:       generator.NowISO8601(),
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update expense: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple expenses updated with the same ID")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "expense not updated")
	}

	return nil
}

// DeleteExpenseByID deletes an expense by its ID if the user has access to the
// book.
//
// It returns true if the expense was deleted, false if the expense does not
// exist or the user does not have access to the book.
func (s *ExpenseService) DeleteExpenseByID(ctx context.Context, userID, expenseID string) error {
	// Get the expense to find the book ID
	expenses, err := s.queries.GetExpenseByID(ctx, expenseID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get expense: %v", err)
	}

	if len(expenses) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple expenses found with the same ID")
	}

	if len(expenses) < 1 {
		return NewServiceError(ErrCodeNotFound, "expense not found or access denied")
	}

	expense := expenses[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: expense.BookID,
		UserID: userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "expense not found or access denied")
	}

	// Proceed to delete the expense
	rows, err := s.queries.DeleteExpenseByID(ctx, expenseID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete expense: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple expenses deleted with the same ID")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "expense not deleted")
	}

	return nil
}

// checkBookCategoryPaymentMethod checks if the user has access to the book,
// category, and payment method.
//
// It also checks if the category and payment method belong to the book.
func (s *ExpenseService) checkBookCategoryPaymentMethod(ctx context.Context, userID, bookID, categoryID, paymentMethodID string) error {
	// Check if the user has access to the book
	canAccessBook, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: bookID,
		UserID: userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccessBook {
		return NewServiceError(ErrCodeUnprocessable, "book not found or access denied")
	}

	// Check if the categories belongs to the book
	categories, err := s.queries.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get category by ID: %v", err)
	}

	if len(categories) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple categories found with the same ID")
	}

	if len(categories) < 1 {
		return NewServiceError(ErrCodeUnprocessable, "category not found or access denied")
	}

	category := categories[0]

	if category.BookID != bookID {
		return NewServiceError(ErrCodeUnprocessable, "category does not belong to the book")
	}

	// Check if the payment method belongs to the book
	paymentMethod, err := s.queries.GetPaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to get payment method by ID: %v", err)
	}

	if len(paymentMethod) > 1 {
		return NewServiceError(ErrCodeInternal, "multiple payment methods found with the same ID")
	}

	if len(paymentMethod) < 1 {
		return NewServiceError(ErrCodeUnprocessable, "payment method not found or access denied")
	}

	pm := paymentMethod[0]

	if pm.BookID != bookID {
		return NewServiceError(ErrCodeUnprocessable, "payment method does not belong to the book")
	}

	return nil
}
