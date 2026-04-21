package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootHandlerServesApiAndSpaFallback(t *testing.T) {
	distDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<html>control room</html>"), 0o600))
	require.NoError(t, os.Mkdir(filepath.Join(distDir, "assets"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(distDir, "assets", "app.js"), []byte("console.log('ok')"), 0o600))

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/ping" {
			_, _ = w.Write([]byte("pong"))
			return
		}
		http.NotFound(w, r)
	})

	handler := newRootHandler(apiHandler, distDir)

	apiResp := httptest.NewRecorder()
	handler.ServeHTTP(apiResp, httptest.NewRequest(http.MethodGet, "/api/ping", nil))
	assert.Equal(t, http.StatusOK, apiResp.Code)
	assert.Equal(t, "pong", apiResp.Body.String())

	assetResp := httptest.NewRecorder()
	handler.ServeHTTP(assetResp, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))
	assert.Equal(t, http.StatusOK, assetResp.Code)
	assert.Contains(t, assetResp.Body.String(), "console.log")

	spaResp := httptest.NewRecorder()
	handler.ServeHTTP(spaResp, httptest.NewRequest(http.MethodGet, "/rooms/room-001", nil))
	assert.Equal(t, http.StatusOK, spaResp.Code)
	assert.Contains(t, spaResp.Body.String(), "control room")
}
