package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/test-go/testify/assert"
)

type healthResponse struct {
	Version string `json:"version"`
}

func TestHealth(t *testing.T) {
	t.Run("returns health check JSON on matching endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(Health("/ping", "app", "v1.0.0-beta.1"))
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should not reach here"))
		})
		r.Get("/other", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("other endpoint"))
		})

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var resp healthResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0-beta.1", resp.Version)
	})

	t.Run("passes through to next handler on non-matching endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(Health("/health", "app", "v1.2.3"))
		r.Get("/other", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("other endpoint"))
		})

		req := httptest.NewRequest(http.MethodGet, "/other", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "other endpoint", rr.Body.String())
	})

	t.Run("handles different endpoint paths", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(Health("/api/health", "app", "v2.0.0"))
		r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should not reach here"))
		})

		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var resp healthResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "v2.0.0", resp.Version)
	})

	t.Run("handles different HTTP methods on health endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(Health("/health", "app", "v1.0.0"))
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should not reach here"))
		})

		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Health middleware should respond regardless of method
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var resp healthResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", resp.Version)
	})
}
