package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/server/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BookID      string `json:"book_id"`
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /categories", h.createCategory)
	mux.HandleFunc("GET /categories", h.getCategoriesByBookID)
	mux.HandleFunc("GET /categories/{id}", h.getCategoryByID)
	mux.HandleFunc("PUT /categories/{id}", h.updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.deleteCategory)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.BookID == "" {
		http.Error(w, "Category name and book ID are required", http.StatusBadRequest)
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

	ok, err := h.categoryService.CreateCategory(ctx, userID, req.BookID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created successfully"))
}

func (h *CategoryHandler) getCategoriesByBookID(w http.ResponseWriter, r *http.Request) {
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

	categories, err := h.categoryService.GetCategoriesByBookID(r.Context(), userID, bookID)
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	if categories == nil {
		http.Error(w, "Book not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	categoryID := r.PathValue("id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
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

	category, err := h.categoryService.GetCategoryByID(r.Context(), userID, categoryID)
	if err != nil {
		http.Error(w, "Failed to get category", http.StatusInternalServerError)
		return
	}

	if category == nil {
		http.Error(w, "Category not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	categoryID := r.PathValue("id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
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

	ok, err := h.categoryService.UpdateCategoryByID(ctx, userID, categoryID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Category not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category updated successfully"))
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Input validation
	categoryID := r.PathValue("id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
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

	ok, err := h.categoryService.DeleteCategoryByID(ctx, userID, categoryID)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Category not found or access denied", http.StatusNotFound)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category deleted successfully"))
}
