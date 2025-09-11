package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/jljl1337/xpense/web"
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	// Get the embedded filesystem rooted at "build"
	buildFS, err := fs.Sub(web.SiteDir, "build/client")
	if err != nil {
		http.Error(w, "Failed to access embedded files", http.StatusInternalServerError)
		return
	}

	// Try to serve the requested file
	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "index.html"
	}

	// Check if file exists
	if _, err := fs.Stat(buildFS, filePath); err == nil {
		// File exists, serve it
		http.FileServer(http.FS(buildFS)).ServeHTTP(w, r)
		return
	}

	// File doesn't exist, serve index.html for SPA routing
	r.URL.Path = "/"
	http.FileServer(http.FS(buildFS)).ServeHTTP(w, r)
}
