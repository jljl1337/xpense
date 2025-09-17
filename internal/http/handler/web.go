package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/jljl1337/xpense/web"
)

type WebHandler struct {
	siteFs fs.FS
}

func NewWebHandler() *WebHandler {
	buildFS, err := fs.Sub(web.SiteDir, "build/client")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}

	return &WebHandler{siteFs: buildFS}
}

func (h *WebHandler) ServeSite(w http.ResponseWriter, r *http.Request) {
	// Try to serve the requested file
	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "index.html"
	}

	// Check if file exists
	if _, err := fs.Stat(h.siteFs, filePath); err == nil {
		// File exists, serve it
		http.FileServer(http.FS(h.siteFs)).ServeHTTP(w, r)
		return
	}

	// File doesn't exist, serve index.html for SPA routing
	r.URL.Path = "/"
	http.FileServer(http.FS(h.siteFs)).ServeHTTP(w, r)
}
