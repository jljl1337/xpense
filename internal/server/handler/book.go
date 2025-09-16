package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/server/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type createUpdateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /books", h.createBook)
	mux.HandleFunc("GET /books", h.getBooks)
	mux.HandleFunc("PUT /books/{id}", h.updateBook)
	mux.HandleFunc("DELETE /books/{id}", h.deleteBook)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	err = h.bookService.CreateBook(ctx, userID, req.Name, req.Description)
	if err != nil {
		slog.Error("Error creating book: " + err.Error())
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created successfully"))
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	// Input validation
	queryValues := r.URL.Query()

	page, err := strconv.ParseInt(queryValues.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(queryValues.Get("page_size"), 10, 64)
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

	books, err := h.bookService.GetBooksByUserID(ctx, userID, page, pageSize)
	if err != nil {
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	ok, err := h.bookService.UpdateBookByID(ctx, userID, bookID, req.Name, req.Description)
	if err != nil {
		slog.Error("Error updating book: " + err.Error())
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated successfully"))
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	ok, err := h.bookService.DeleteBookByID(ctx, userID, bookID)
	if err != nil {
		slog.Error("Error deleting book: " + err.Error())
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted successfully"))
}
