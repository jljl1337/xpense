package service

import (
	"context"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

// CreatePaymentMethod creates a new payment method if the user has access to the book.
func (s *EndpointService) CreatePaymentMethod(ctx context.Context, userID, bookID, name, description string) error {
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
		return NewServiceError(ErrCodeUnprocessable, "book not found or access denied")
	}

	currentTime := generator.NowISO8601()

	_, err = queries.CreatePaymentMethod(ctx, repository.CreatePaymentMethodParams{
		ID:          generator.NewULID(),
		BookID:      bookID,
		Name:        name,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to create payment method: %v", err)
	}

	return nil
}

// GetPaymentMethodsByBookID retrieves all payment methods for a specific book.
//
// It returns an empty slice if no payment methods are found in the book.
func (s *EndpointService) GetPaymentMethodsByBookID(ctx context.Context, userID, bookID string) ([]repository.PaymentMethod, error) {
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
		return nil, NewServiceError(ErrCodeUnprocessable, "book not found or access denied")
	}

	paymentMethods, err := queries.GetPaymentMethodsByBookID(ctx, bookID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get payment methods by book ID: %v", err)
	}

	return paymentMethods, nil
}

// GetPaymentMethodByID retrieves a payment method by its ID if the user has access to the book.
func (s *EndpointService) GetPaymentMethodByID(ctx context.Context, userID, paymentMethodID string) (*repository.PaymentMethod, error) {
	queries := repository.New(s.db)

	// Get the payment method to find the book ID
	paymentMethods, err := queries.GetPaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to get payment method: %v", err)
	}

	if len(paymentMethods) > 1 {
		return nil, NewServiceError(ErrCodeInternal, "multiple payment methods found with the same ID")
	}

	if len(paymentMethods) < 1 {
		return nil, NewServiceError(ErrCodeNotFound, "payment method not found or access denied")
	}

	paymentMethod := paymentMethods[0]

	// Check if the user has access to the book
	canAccess, err := queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: paymentMethod.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, NewServiceErrorf(ErrCodeInternal, "failed to check book access: %v", err)
	}

	if !canAccess {
		return nil, NewServiceError(ErrCodeNotFound, "book not found or access denied")
	}

	return &paymentMethod, nil
}

// UpdatePaymentMethodByID updates a payment method if the user has access to the book.
func (s *EndpointService) UpdatePaymentMethodByID(ctx context.Context, userID, paymentMethodID, name, description string) error {
	queries := repository.New(s.db)

	// Check if the user has access to the payment method
	canAccess, err := queries.CheckPaymentMethodAccess(ctx, repository.CheckPaymentMethodAccessParams{
		PaymentMethodID: paymentMethodID,
		UserID:          userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check payment method access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "payment method not found or access denied")
	}

	rows, err := queries.UpdatePaymentMethodByID(ctx, repository.UpdatePaymentMethodByIDParams{
		ID:          paymentMethodID,
		Name:        name,
		Description: description,
		UpdatedAt:   generator.NowISO8601(),
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to update payment method: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple payment methods updated, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no payment method updated")
	}

	return nil
}

// DeletePaymentMethodByID deletes a payment method if the user has access to the book.
func (s *EndpointService) DeletePaymentMethodByID(ctx context.Context, userID, paymentMethodID string) error {
	queries := repository.New(s.db)

	// Check if the user has access to the payment method
	canAccess, err := queries.CheckPaymentMethodAccess(ctx, repository.CheckPaymentMethodAccessParams{
		PaymentMethodID: paymentMethodID,
		UserID:          userID,
	})
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to check payment method access: %v", err)
	}

	if !canAccess {
		return NewServiceError(ErrCodeNotFound, "payment method not found or access denied")
	}

	rows, err := queries.DeletePaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return NewServiceErrorf(ErrCodeInternal, "failed to delete payment method: %v", err)
	}

	if rows > 1 {
		return NewServiceError(ErrCodeInternal, "multiple payment methods deleted, data integrity issue")
	}

	if rows < 1 {
		return NewServiceError(ErrCodeInternal, "no payment method deleted")
	}

	return nil
}
