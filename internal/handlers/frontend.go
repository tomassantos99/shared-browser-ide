package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	staticDir, err := filepath.Abs("./static")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedPath := filepath.Join(staticDir, r.URL.Path)

	// Serve static file if it exists (js/css or index.html)
	if info, err := os.Stat(requestedPath); err == nil && !info.IsDir() {
		http.ServeFile(w, r, requestedPath)
		return
	}

	// Default
	http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
}
