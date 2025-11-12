package handler

import (
	"net/http"

	"github.com/jljl1337/xpense/internal/service"
)

type EndpointHandler struct {
	service *service.EndpointService
}

func NewEndpointHandler(service *service.EndpointService) *EndpointHandler {
	return &EndpointHandler{
		service: service,
	}
}

func (h *EndpointHandler) RegisterRoutes(mux *http.ServeMux) {
	h.registerAuthRoutes(mux)
	h.registerUserRoutes(mux)
	h.registerBookRoutes(mux)
	h.registerCategoryRoutes(mux)
	h.registerPaymentMethodRoutes(mux)
	h.registerExpenseRoutes(mux)
	h.registerHealthCheckRoutes(mux)
	h.registerVersionRoutes(mux)
}
