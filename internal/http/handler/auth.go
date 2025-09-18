package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type signUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginCSRFTokenResponse struct {
	CSRFToken string `json:"csrf_token"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/sign-up", h.signUp)
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/logout", h.logout)
	mux.HandleFunc("POST /auth/logout-all", h.logoutAll)
	mux.HandleFunc("GET /auth/csrf-token", h.csrfToken)
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req signUpLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Process the request
	if err := h.authService.SignUp(r.Context(), req.Email, req.Password); err != nil {
		slog.Error("Error signing up user: " + err.Error())
		http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User signed up successfully"))
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req signUpLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Process the request
	sessionToken, CSRFToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("Error logging in user: " + err.Error())
		http.Error(w, "Failed to log in user", http.StatusInternalServerError)
		return
	}

	if sessionToken == "" && CSRFToken == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Respond to the client
	http.SetCookie(w, &http.Cookie{
		Name:     env.SessionCookieName,
		Value:    sessionToken,
		HttpOnly: env.SessionCookieHttpOnly,
		Secure:   env.SessionCookieSecure,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	// Input validation
	sessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request
	if err := h.authService.Logout(r.Context(), sessionToken.Value); err != nil {
		slog.Error("Error logging out user: " + err.Error())
		http.Error(w, "Failed to log out user", http.StatusInternalServerError)
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
	w.Write([]byte("User logged out successfully"))
}

func (h *AuthHandler) logoutAll(w http.ResponseWriter, r *http.Request) {
	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	if err := h.authService.LogoutAllSessions(r.Context(), userID); err != nil {
		slog.Error("Error logging out user from all sessions: " + err.Error())
		http.Error(w, "Failed to log out user from all sessions", http.StatusInternalServerError)
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
	w.Write([]byte("User logged out from all sessions successfully"))
}

func (h *AuthHandler) csrfToken(w http.ResponseWriter, r *http.Request) {
	// Input validation
	sessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request
	CSRFToken, err := h.authService.CSRFToken(r.Context(), sessionToken.Value)
	if err != nil {
		slog.Error("Error getting CSRF token: " + err.Error())
		http.Error(w, "Failed to get CSRF token", http.StatusInternalServerError)
		return
	}

	if CSRFToken == "" {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}
