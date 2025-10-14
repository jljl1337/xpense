package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
)

type versionResponse struct {
	Version string `json:"version"`
}

type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

func (h *VersionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/version", h.version)
}

func (h *VersionHandler) version(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionResponse{Version: env.Version})
}
