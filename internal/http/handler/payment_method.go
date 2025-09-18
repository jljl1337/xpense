package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type createPaymentMethodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BookID      string `json:"book_id"`
}

type updatePaymentMethodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PaymentMethodHandler struct {
	paymentMethodService *service.PaymentMethodService
}

func NewPaymentMethodHandler(paymentMethodService *service.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentMethodService: paymentMethodService,
	}
}

func (h *PaymentMethodHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /payment-methods", h.createPaymentMethod)
	mux.HandleFunc("GET /payment-methods", h.getPaymentMethodsByBookID)
	mux.HandleFunc("GET /payment-methods/{id}", h.getPaymentMethodByID)
	mux.HandleFunc("PUT /payment-methods/{id}", h.updatePaymentMethod)
	mux.HandleFunc("DELETE /payment-methods/{id}", h.deletePaymentMethod)
}

func (h *PaymentMethodHandler) createPaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req createPaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.BookID == "" {
		http.Error(w, "Payment method name and book ID are required", http.StatusBadRequest)
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

	ok, err := h.paymentMethodService.CreatePaymentMethod(ctx, userID, req.BookID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create payment method", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Payment method created successfully"))
}

func (h *PaymentMethodHandler) getPaymentMethodsByBookID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
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

	paymentMethods, err := h.paymentMethodService.GetPaymentMethodsByBookID(r.Context(), userID, bookID)
	if err != nil {
		http.Error(w, "Failed to get payment methods", http.StatusInternalServerError)
		return
	}

	if paymentMethods == nil {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentMethods)
}

func (h *PaymentMethodHandler) getPaymentMethodByID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	paymentMethodID := r.PathValue("id")
	if paymentMethodID == "" {
		http.Error(w, "Payment method ID is required", http.StatusBadRequest)
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

	paymentMethod, err := h.paymentMethodService.GetPaymentMethodByID(r.Context(), userID, paymentMethodID)
	if err != nil {
		http.Error(w, "Failed to get payment method", http.StatusInternalServerError)
		return
	}

	if paymentMethod == nil {
		http.Error(w, "Payment method not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentMethod)
}

func (h *PaymentMethodHandler) updatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req updatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Payment method name is required", http.StatusBadRequest)
		return
	}

	paymentMethodID := r.PathValue("id")
	if paymentMethodID == "" {
		http.Error(w, "Payment method ID is required", http.StatusBadRequest)
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

	ok, err := h.paymentMethodService.UpdatePaymentMethodByID(ctx, userID, paymentMethodID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to update payment method", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Payment method not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment method updated successfully"))
}

func (h *PaymentMethodHandler) deletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Input validation
	paymentMethodID := r.PathValue("id")
	if paymentMethodID == "" {
		http.Error(w, "Payment method ID is required", http.StatusBadRequest)
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

	ok, err := h.paymentMethodService.DeletePaymentMethodByID(ctx, userID, paymentMethodID)
	if err != nil {
		http.Error(w, "Failed to delete payment method", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Payment method not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment method deleted successfully"))
}
