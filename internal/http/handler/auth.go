package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type signUpSignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInPreSessionCSRFTokenResponse struct {
	CSRFToken string `json:"csrfToken"`
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
	mux.HandleFunc("POST /auth/pre-session", h.preSession)
	mux.HandleFunc("POST /auth/sign-in", h.signIn)
	mux.HandleFunc("POST /auth/sign-out", h.signOut)
	mux.HandleFunc("POST /auth/sign-out-all", h.signOutAll)
	mux.HandleFunc("GET /auth/csrf-token", h.csrfToken)
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	// Input validation
	var req signUpSignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Process the request
	if err := h.authService.SignUp(r.Context(), req.Username, req.Password); err != nil {
		slog.Error("Error signing up user: " + err.Error())
		http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User signed up successfully"))
}

func (h *AuthHandler) preSession(w http.ResponseWriter, r *http.Request) {
	// Process the request
	sessionToken, CSRFToken, err := h.authService.GetPreSession(r.Context())
	if err != nil {
		slog.Error("Error getting pre-session: " + err.Error())
		http.Error(w, "Failed to get pre-session", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewActiveSessionCookie(sessionToken))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	// Input validation
	preSessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	preSessionCSRFToken := r.Header.Get("X-CSRF-Token")
	if preSessionCSRFToken == "" {
		http.Error(w, "CSRF token is required", http.StatusUnauthorized)
		return
	}

	var req signUpSignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Process the request
	sessionToken, CSRFToken, err := h.authService.SignIn(r.Context(), preSessionToken.Value, preSessionCSRFToken, req.Username, req.Password)
	if err != nil {
		slog.Error("Error signing in user: " + err.Error())
		http.Error(w, "Failed to sign in user", http.StatusInternalServerError)
		return
	}

	if sessionToken == "" && CSRFToken == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewActiveSessionCookie(sessionToken))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	// Input validation
	sessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request
	if err := h.authService.SignOut(r.Context(), sessionToken.Value); err != nil {
		slog.Error("Error signing out user: " + err.Error())
		http.Error(w, "Failed to sign out user", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewExpiredSessionCookie())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}

func (h *AuthHandler) signOutAll(w http.ResponseWriter, r *http.Request) {
	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Failed to get the current user", http.StatusInternalServerError)
		return
	}

	if err := h.authService.SignOutAllSession(r.Context(), userID); err != nil {
		slog.Error("Error signing out user from all sessions: " + err.Error())
		http.Error(w, "Failed to sign out user from all sessions", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewExpiredSessionCookie())

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
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}
