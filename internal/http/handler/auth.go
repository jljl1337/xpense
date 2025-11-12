package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/common"
	"github.com/jljl1337/xpense/internal/http/middleware"
)

type signUpSignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInPreSessionCSRFTokenResponse struct {
	CSRFToken string `json:"csrfToken"`
}

func (h *EndpointHandler) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/sign-up", h.signUp)
	mux.HandleFunc("POST /auth/pre-session", h.preSession)
	mux.HandleFunc("POST /auth/sign-in", h.signIn)
	mux.HandleFunc("POST /auth/sign-out", h.signOut)
	mux.HandleFunc("POST /auth/sign-out-all", h.signOutAll)
	mux.HandleFunc("GET /auth/csrf-token", h.csrfToken)
}

func (h *EndpointHandler) signUp(w http.ResponseWriter, r *http.Request) {
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
	if err := h.service.SignUp(r.Context(), req.Username, req.Password); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User signed up successfully"))
}

func (h *EndpointHandler) preSession(w http.ResponseWriter, r *http.Request) {
	// Process the request
	sessionToken, CSRFToken, err := h.service.GetPreSession(r.Context())
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewActiveSessionCookie(sessionToken))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *EndpointHandler) signIn(w http.ResponseWriter, r *http.Request) {
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
	sessionToken, CSRFToken, err := h.service.SignIn(r.Context(), preSessionToken.Value, preSessionCSRFToken, req.Username, req.Password)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewActiveSessionCookie(sessionToken))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}

func (h *EndpointHandler) signOut(w http.ResponseWriter, r *http.Request) {
	// Input validation
	sessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request
	if err := h.service.SignOut(r.Context(), sessionToken.Value); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewExpiredSessionCookie())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}

func (h *EndpointHandler) signOutAll(w http.ResponseWriter, r *http.Request) {
	// Process the request
	ctx := r.Context()
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Error getting user ID from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.service.SignOutAllSession(r.Context(), userID); err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	http.SetCookie(w, NewExpiredSessionCookie())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out from all sessions successfully"))
}

func (h *EndpointHandler) csrfToken(w http.ResponseWriter, r *http.Request) {
	// Input validation
	sessionToken, err := r.Cookie(env.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request
	CSRFToken, err := h.service.CSRFToken(r.Context(), sessionToken.Value)
	if err != nil {
		common.WriteErrorResponse(w, err)
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signInPreSessionCSRFTokenResponse{
		CSRFToken: CSRFToken,
	})
}
