package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/service"
)

type signUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
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
	mux.HandleFunc("POST /auth/sign-up", h.signUpHandler)
	mux.HandleFunc("POST /auth/login", h.loginHandler)
	mux.HandleFunc("POST /auth/logout", h.logoutHandler)
}

func (h *AuthHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var req signUpLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if err := h.authService.SignUp(req.Email, req.Password); err != nil {
		slog.Error("Error signing up user: " + err.Error())
		http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User signed up successfully"))
}

func (h *AuthHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req signUpLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	sessionToken, CSRFToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		slog.Error("Error logging in user: " + err.Error())
		http.Error(w, "Failed to log in user", http.StatusInternalServerError)
		return
	}

	if sessionToken == "" && CSRFToken == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *AuthHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.authService.Logout(sessionToken.Value); err != nil {
		slog.Error("Error logging out user: " + err.Error())
		http.Error(w, "Failed to log out user", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}
