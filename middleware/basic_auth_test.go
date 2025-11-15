package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/test-go/testify/assert"

	"github.com/0xsequence/go-libs/config"
)

func TestBasicAuth(t *testing.T) {
	t.Run("valid creds", func(t *testing.T) {
		creds := config.BasicAuth{
			Username: "testuser",
			Password: "testpass",
		}

		r := chi.NewRouter()
		r.Use(BasicAuth(creds))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("protected content"))
		})

		t.Run("returns 401 without credentials", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
			assert.Contains(t, rr.Header().Get("WWW-Authenticate"), "Basic")
		})

		t.Run("returns 401 with invalid credentials", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.SetBasicAuth("wronguser", "wrongpass")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})

		t.Run("returns 200 with valid credentials", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.SetBasicAuth("testuser", "testpass")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, "protected content", rr.Body.String())
		})
	})

	t.Run("missing creds", func(t *testing.T) {
		creds := config.BasicAuth{
			Username: "",
			Password: "",
		}

		r := chi.NewRouter()
		r.Use(BasicAuth(creds))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("protected content"))
		})

		t.Run("returns 404 when credentials are empty", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNotFound, rr.Code)
			assert.Equal(t, "Not Found\n", rr.Body.String())
		})
	})

	t.Run("with partial credentials configured", func(t *testing.T) {
		creds := config.BasicAuth{
			Username: "testuser",
			Password: "",
		}

		r := chi.NewRouter()
		r.Use(BasicAuth(creds))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("protected content"))
		})

		t.Run("returns 404 when password is empty", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNotFound, rr.Code)
			assert.Equal(t, "Not Found\n", rr.Body.String())
		})
	})
}
