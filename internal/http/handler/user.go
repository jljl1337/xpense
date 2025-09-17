package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type getCurrentUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/me", h.getCurrentUser)
	mux.HandleFunc("DELETE /users/me", h.deleteCurrentUser)
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if user == nil || err != nil {
		slog.Error("Failed to get user with ID: " + userID)
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	response := getCurrentUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) deleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to delete the current user", http.StatusInternalServerError)
		return
	}

	if err := h.userService.DeleteUserByID(r.Context(), userID); err != nil {
		slog.Error("Error deleting user with ID: " + userID + " - " + err.Error())
		http.Error(w, "Failed to delete the current user", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
