package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func newRootHandler(apiHandler http.Handler, webDist string) http.Handler {
	if webDist == "" {
		return apiHandler
	}

	fileServer := http.FileServer(http.Dir(webDist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws/") {
			apiHandler.ServeHTTP(w, r)
			return
		}

		cleanPath := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if cleanPath == "." || cleanPath == "" {
			http.ServeFile(w, r, filepath.Join(webDist, "index.html"))
			return
		}

		fullPath := filepath.Join(webDist, filepath.FromSlash(cleanPath))
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(webDist, "index.html"))
	})
}
