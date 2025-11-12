package repository

import (
	"database/sql"
)

type Book struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"userID" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	UpdatedAt   string `json:"updatedAt" db:"updated_at"`
}

type Category struct {
	ID          string `json:"id" db:"id"`
	BookID      string `json:"bookID" db:"book_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	UpdatedAt   string `json:"updatedAt" db:"updated_at"`
}

type Expense struct {
	ID              string  `json:"id" db:"id"`
	BookID          string  `json:"bookID" db:"book_id"`
	CategoryID      string  `json:"categoryID" db:"category_id"`
	PaymentMethodID string  `json:"paymentMethodID" db:"payment_method_id"`
	Date            string  `json:"date" db:"date"`
	Amount          float64 `json:"amount" db:"amount"`
	Remark          string  `json:"remark" db:"remark"`
	CreatedAt       string  `json:"createdAt" db:"created_at"`
	UpdatedAt       string  `json:"updatedAt" db:"updated_at"`
}

type PaymentMethod struct {
	ID          string `json:"id" db:"id"`
	BookID      string `json:"bookID" db:"book_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	UpdatedAt   string `json:"updatedAt" db:"updated_at"`
}

type Session struct {
	ID        string         `json:"id" db:"id"`
	UserID    sql.NullString `json:"userID" db:"user_id"`
	Token     string         `json:"token" db:"token"`
	CsrfToken string         `json:"csrfToken" db:"csrf_token"`
	ExpiresAt string         `json:"expiresAt" db:"expires_at"`
	CreatedAt string         `json:"createdAt" db:"created_at"`
	UpdatedAt string         `json:"updatedAt" db:"updated_at"`
}

type User struct {
	ID           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"passwordHash" db:"password_hash"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
	UpdatedAt    string `json:"updatedAt" db:"updated_at"`
}
