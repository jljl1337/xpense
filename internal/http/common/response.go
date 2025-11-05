package common

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/service"
)

func WriteErrorResponse(w http.ResponseWriter, err error) {
	var serviceErr *service.ServiceError
	if errors.As(err, &serviceErr) {
		httpStatus := mapServiceErrorToHTTPStatus(serviceErr)

		if httpStatus == http.StatusInternalServerError {
			slog.Error("Internal server error: " + serviceErr.Error())
			http.Error(w, "Internal server error", httpStatus)
			return
		}

		http.Error(w, serviceErr.Message, httpStatus)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func mapServiceErrorToHTTPStatus(err *service.ServiceError) int {
	switch err.Code {
	case service.ErrCodeBadRequest:
		return http.StatusBadRequest
	case service.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case service.ErrCodeForbidden:
		return http.StatusForbidden
	case service.ErrCodeNotFound:
		return http.StatusNotFound
	case service.ErrCodeConflict:
		return http.StatusConflict
	case service.ErrCodeUnprocessable:
		return http.StatusUnprocessableEntity
	case service.ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
