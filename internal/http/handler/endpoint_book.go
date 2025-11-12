package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/common"
	"github.com/jljl1337/xpense/internal/http/middleware"
)

type createUpdateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type getBooksCountResponse struct {
	Count int64 `json:"count"`
}

func (h *EndpointHandler) registerBookRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /books", h.createBook)
	mux.HandleFunc("GET /books/count", h.getBooksCount)
	mux.HandleFunc("GET /books", h.getBooks)
	mux.HandleFunc("GET /books/{id}", h.getBook)
	mux.HandleFunc("PUT /books/{id}", h.updateBook)
	mux.HandleFunc("DELETE /books/{id}", h.deleteBook)
}

func (h *EndpointHandler) createBook(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req createUpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Book name is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.CreateBook(ctx, userID, req.Name, req.Description)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created successfully"))
}

func (h *EndpointHandler) getBooksCount(w http.ResponseWriter, r *http.Request) {
	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	count, err := h.service.GetBooksCountByUserID(ctx, userID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getBooksCountResponse{Count: count})
}

func (h *EndpointHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	// Input validation
	queryValues := r.URL.Query()

	page, err := strconv.ParseInt(queryValues.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(queryValues.Get("page-size"), 10, 64)
	if err != nil || pageSize < 1 || pageSize > env.PageSizeMax {
		pageSize = env.PageSizeDefault
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	books, err := h.service.GetBooksByUserID(ctx, userID, page, pageSize)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *EndpointHandler) getBook(w http.ResponseWriter, r *http.Request) {
	// Input validation
	bookID := r.PathValue("id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	book, err := h.service.GetBookByID(ctx, userID, bookID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *EndpointHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req createUpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	bookID := r.PathValue("id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Book name is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.UpdateBookByID(ctx, userID, bookID, req.Name, req.Description)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated successfully"))
}

func (h *EndpointHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	// Input validation
	bookID := r.PathValue("id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteBookByID(ctx, userID, bookID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted successfully"))
}
