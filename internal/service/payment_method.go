package service

import (
	"context"
	"errors"
	"time"

	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

type PaymentMethodService struct {
	queries *repository.Queries
}

func NewPaymentMethodService(queries *repository.Queries) *PaymentMethodService {
	return &PaymentMethodService{
		queries: queries,
	}
}

// CreatePaymentMethod creates a new payment method if the user has access to the book.
//
// It returns true if the payment method was created, false if the user has no access to the book
func (s *PaymentMethodService) CreatePaymentMethod(ctx context.Context, userID, bookID, name, description string) (bool, error) {
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

	_, err = s.queries.CreatePaymentMethod(ctx, repository.CreatePaymentMethodParams{
		ID:          generator.NewKSUID(),
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

// GetPaymentMethodsByBookID retrieves all payment methods for a specific book.
//
// It returns an empty slice if no payment methods are found in the book.
//
// It returns nil if the user has no access to the book.
func (s *PaymentMethodService) GetPaymentMethodsByBookID(ctx context.Context, userID, bookID string) ([]repository.PaymentMethod, error) {
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

	paymentMethods, err := s.queries.GetPaymentMethodsByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return paymentMethods, nil
}

// GetPaymentMethodByID retrieves a payment method by its ID if the user has access to the book.
//
// It returns nil if the payment method does not exist or the user does not have
// access to the book.
func (s *PaymentMethodService) GetPaymentMethodByID(ctx context.Context, userID, paymentMethodID string) (*repository.PaymentMethod, error) {
	// Get the payment method to find the book ID
	paymentMethods, err := s.queries.GetPaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return nil, err
	}

	if len(paymentMethods) > 1 {
		return nil, errors.New("multiple payment methods found with the same ID")
	}

	if len(paymentMethods) < 1 {
		return nil, nil
	}

	paymentMethod := paymentMethods[0]

	// Check if the user has access to the book
	canAccess, err := s.queries.CheckBookAccess(ctx, repository.CheckBookAccessParams{
		BookID: paymentMethod.BookID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if !canAccess {
		return nil, nil
	}

	return &paymentMethod, nil
}

// UpdatePaymentMethodByID updates a payment method if the user has access to the book.
//
// It returns true if the payment method was updated, false if the user has no access
// to the payment method or the payment method does not exist
func (s *PaymentMethodService) UpdatePaymentMethodByID(ctx context.Context, userID, paymentMethodID, name, description string) (bool, error) {
	// Check if the user has access to the payment method
	canAccess, err := s.queries.CheckPaymentMethodAccess(ctx, repository.CheckPaymentMethodAccessParams{
		PaymentMethodID: paymentMethodID,
		UserID:          userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	currentTime := time.Now().UnixMilli()

	rows, err := s.queries.UpdatePaymentMethodByID(ctx, repository.UpdatePaymentMethodByIDParams{
		ID:          paymentMethodID,
		Name:        name,
		Description: description,
		UpdatedAt:   currentTime,
	})
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple payment methods updated, data integrity issue")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}

// DeletePaymentMethodByID deletes a payment method if the user has access to the book.
//
// It returns true if the payment method was deleted, false if the user has no access
// to the payment method or the payment method does not exist
func (s *PaymentMethodService) DeletePaymentMethodByID(ctx context.Context, userID, paymentMethodID string) (bool, error) {
	// Check if the user has access to the payment method
	canAccess, err := s.queries.CheckPaymentMethodAccess(ctx, repository.CheckPaymentMethodAccessParams{
		PaymentMethodID: paymentMethodID,
		UserID:          userID,
	})
	if err != nil {
		return false, err
	}

	if !canAccess {
		return false, nil
	}

	rows, err := s.queries.DeletePaymentMethodByID(ctx, paymentMethodID)
	if err != nil {
		return false, err
	}

	if rows > 1 {
		return false, errors.New("multiple payment methods deleted, data integrity issue")
	}

	if rows < 1 {
		return false, nil
	}

	return true, nil
}
