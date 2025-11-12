package middleware

import "github.com/jljl1337/xpense/internal/service"

// MiddlewareProvider contains all middleware functions
type MiddlewareProvider struct {
	service *service.MiddlewareService
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(service *service.MiddlewareService) *MiddlewareProvider {
	return &MiddlewareProvider{
		service: service,
	}
}
