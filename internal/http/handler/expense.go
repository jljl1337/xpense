package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type createExpenseRequest struct {
	BookID          string  `json:"bookID"`
	CategoryID      string  `json:"categoryID"`
	PaymentMethodID string  `json:"paymentMethodID"`
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	Remark          string  `json:"remark"`
}

type updateExpenseRequest struct {
	CategoryID      string  `json:"categoryID"`
	PaymentMethodID string  `json:"paymentMethodID"`
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	Remark          string  `json:"remark"`
}

type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

func (h *ExpenseHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /expenses", h.createExpense)
	mux.HandleFunc("GET /expenses", h.getExpensesByBookID)
	mux.HandleFunc("GET /expenses/{id}", h.getExpenseByID)
	mux.HandleFunc("PUT /expenses/{id}", h.updateExpense)
	mux.HandleFunc("DELETE /expenses/{id}", h.deleteExpense)
}

func (h *ExpenseHandler) createExpense(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.BookID == "" || req.CategoryID == "" || req.PaymentMethodID == "" {
		http.Error(w, "Book ID, category ID and payment method ID are required", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		http.Error(w, "Date must be a valid YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	created, err := h.expenseService.CreateExpense(ctx, userID, req.BookID, req.CategoryID, req.PaymentMethodID, req.Date, req.Amount, req.Remark)
	if err != nil {
		slog.Error("Error creating expense: " + err.Error())
		http.Error(w, "Failed to create expense", http.StatusInternalServerError)
		return
	}

	if !created {
		http.Error(w, "Failed to create expense. Make sure you have access to the book, category and payment method, and that the category and payment method belong to the book.", http.StatusBadRequest)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Expense created successfully"))
}

func (h *ExpenseHandler) getExpensesByBookID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	bookID := r.URL.Query().Get("book-id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil || pageSize < 1 || pageSize > env.PageSizeMax {
		pageSize = env.PageSizeDefault
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	expenses, err := h.expenseService.GetExpensesByBookID(r.Context(), userID, bookID, page, pageSize)
	if err != nil {
		slog.Error("Error getting expenses: " + err.Error())
		http.Error(w, "Failed to get expenses", http.StatusInternalServerError)
		return
	}

	if expenses == nil {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpenseHandler) getExpenseByID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	expense, err := h.expenseService.GetExpenseByID(r.Context(), userID, expenseID)
	if err != nil {
		slog.Error("Error getting expense: " + err.Error())
		http.Error(w, "Failed to get expense", http.StatusInternalServerError)
		return
	}

	if expense == nil {
		http.Error(w, "Expense not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) updateExpense(w http.ResponseWriter, r *http.Request) {
	// Input validation
	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	var req updateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.CategoryID == "" || req.PaymentMethodID == "" {
		http.Error(w, "Category ID and payment method ID are required", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		http.Error(w, "Date must be a valid YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	updated, err := h.expenseService.UpdateExpense(ctx, userID, expenseID, req.CategoryID, req.PaymentMethodID, req.Date, req.Amount, req.Remark)
	if err != nil {
		slog.Error("Error updating expense: " + err.Error())
		http.Error(w, "Failed to update expense", http.StatusInternalServerError)
		return
	}

	if !updated {
		http.Error(w, "Failed to update expense. Make sure you have access to the book, category and payment method, and that the category and payment method belong to the book.", http.StatusBadRequest)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expense updated successfully"))
}

func (h *ExpenseHandler) deleteExpense(w http.ResponseWriter, r *http.Request) {
	// Input validation
	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	deleted, err := h.expenseService.DeleteExpenseByID(r.Context(), userID, expenseID)
	if err != nil {
		slog.Error("Error deleting expense: " + err.Error())
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}

	if !deleted {
		http.Error(w, "Expense not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expense deleted successfully"))
}
