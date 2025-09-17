package middleware

import "github.com/jljl1337/xpense/internal/service"

// MiddlewareProvider contains all middleware functions
type MiddlewareProvider struct {
	authService *service.AuthService
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(authService *service.AuthService) *MiddlewareProvider {
	return &MiddlewareProvider{
		authService: authService,
	}
}
