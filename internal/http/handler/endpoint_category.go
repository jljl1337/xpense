package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/http/common"
	"github.com/jljl1337/xpense/internal/http/middleware"
)

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BookID      string `json:"bookID"`
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *EndpointHandler) registerCategoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /categories", h.createCategory)
	mux.HandleFunc("GET /categories", h.getCategoriesByBookID)
	mux.HandleFunc("GET /categories/{id}", h.getCategoryByID)
	mux.HandleFunc("PUT /categories/{id}", h.updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.deleteCategory)
}

func (h *EndpointHandler) createCategory(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.CreateCategory(ctx, userID, req.BookID, req.Name, req.Description)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created successfully"))
}

func (h *EndpointHandler) getCategoriesByBookID(w http.ResponseWriter, r *http.Request) {
	// Input validation
	bookID := r.URL.Query().Get("book-id")
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

	categories, err := h.service.GetCategoriesByBookID(r.Context(), userID, bookID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *EndpointHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	category, err := h.service.GetCategoryByID(r.Context(), userID, categoryID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *EndpointHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.UpdateCategoryByID(ctx, userID, categoryID, req.Name, req.Description)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category updated successfully"))
}

func (h *EndpointHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteCategoryByID(ctx, userID, categoryID)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category deleted successfully"))
}
