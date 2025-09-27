package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type usernameExistResponse struct {
	Exists bool `json:"exists"`
}

type getCurrentUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
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
	mux.HandleFunc("GET /users/exists", h.getUsernameExist)
	mux.HandleFunc("GET /users/me", h.getCurrentUser)
	mux.HandleFunc("DELETE /users/me", h.deleteCurrentUser)
}

func (h *UserHandler) getUsernameExist(w http.ResponseWriter, r *http.Request) {
	// Input validation
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Process the request
	exists, err := h.userService.UserExistsByUsername(r.Context(), username)
	if err != nil {
		slog.Error("Error checking username existence: " + err.Error())
		http.Error(w, "Failed to check username existence", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	response := usernameExistResponse{Exists: exists}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		Username:  user.Username,
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
	http.SetCookie(w, &http.Cookie{
		Name:     env.SessionCookieName,
		Value:    "",
		HttpOnly: env.SessionCookieHttpOnly,
		Secure:   env.SessionCookieSecure,
		Path:     "/",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
