package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/http/common"
	"github.com/jljl1337/xpense/internal/http/middleware"
)

type usernameExistResponse struct {
	Exists bool `json:"exists"`
}

type getCurrentUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

func (h *EndpointHandler) registerUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/exists", h.getUsernameExist)
	mux.HandleFunc("GET /users/me", h.getCurrentUser)
	mux.HandleFunc("PATCH /users/me/username", h.updateUsername)
	mux.HandleFunc("PATCH /users/me/password", h.updatePassword)
	mux.HandleFunc("DELETE /users/me", h.deleteCurrentUser)
}

func (h *EndpointHandler) getUsernameExist(w http.ResponseWriter, r *http.Request) {
	// Input validation
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Process the request
	exists, err := h.service.UserExistsByUsername(r.Context(), username)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	response := usernameExistResponse{Exists: exists}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *EndpointHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		common.WriteErrorResponse(w, err)
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

func (h *EndpointHandler) updateUsername(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req struct {
		NewUsername string `json:"newUsername"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.NewUsername == "" {
		http.Error(w, "New username is required", http.StatusBadRequest)
		return
	}

	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdateUsernameByID(r.Context(), userID, req.NewUsername); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Username updated successfully"))
}

func (h *EndpointHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.NewPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdatePasswordByID(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}

func (h *EndpointHandler) deleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Process the request
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.service.DeleteUserByID(r.Context(), userID); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewExpiredSessionCookie())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
