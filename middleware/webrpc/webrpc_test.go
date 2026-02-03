package webrpc_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/metrics"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/test-go/testify/assert"

	"github.com/0xsequence/go-libs/middleware/webrpc"
)

func TestWebrpcTelemetry(t *testing.T) {
	t.Run("no origin label", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(metrics.Collector(metrics.CollectorOpts{
			Host:  false,
			Proto: true,
			Skip: func(r *http.Request) bool {
				return r.Method != "OPTIONS"
			},
		}))
		r.Use(webrpc.Telemetry(webrpc.Opts{}))
		r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
		})
		r.Handle("/metrics", metrics.Handler())

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://disabled-test-telemetry.example/path?x=y#z")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, 422, rr.Code)

		mfs := scrapeMetrics(t, r)
		mf := mfs["webrpc_requests_total"]
		assert.NotNil(t, mf)
		assert.True(t, metricHasLabels(mf, map[string]string{"status": "422", "origin": ""}))
		assert.False(t, metricHasLabels(mf, map[string]string{"status": "422", "origin": "disabled-test-telemetry.example"}))
	})

	t.Run("origin label with host only", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(metrics.Collector(metrics.CollectorOpts{
			Host:  false,
			Proto: true,
			Skip: func(r *http.Request) bool {
				return r.Method != "OPTIONS"
			},
		}))
		r.Use(webrpc.Telemetry(webrpc.Opts{Origin: true}))
		r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
		})
		r.Handle("/metrics", metrics.Handler())

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://enabled-test-telemetry.example/some/path?x=y#z")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, 422, rr.Code)

		mfs := scrapeMetrics(t, r)
		mf := mfs["webrpc_requests_total"]
		assert.NotNil(t, mf)
		assert.True(t, metricHasLabels(mf, map[string]string{"status": "422", "origin": "enabled-test-telemetry.example"}))
		assert.False(t, metricHasLabels(mf, map[string]string{"status": "422", "origin": "enabled-test-telemetry.example/some/path?x=y#z"}))
	})

	t.Run("OPTIONS preflight can be skipped", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(metrics.Collector(metrics.CollectorOpts{
			Host:  false,
			Proto: true,
			Skip: func(r *http.Request) bool {
				return r.Method != "OPTIONS"
			},
		}))
		r.Use(webrpc.Telemetry(webrpc.Opts{
			Origin: true,
			Skip: func(r *http.Request) bool {
				// Typical CORS preflight signal; avoids dropping legitimate OPTIONS.
				return r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != ""
			},
		}))
		r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
		})
		r.Handle("/metrics", metrics.Handler())

		req := httptest.NewRequest(http.MethodOptions, "/ok", nil)
		req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		req.Header.Set("Origin", "https://no-header.example/")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Handler is not invoked because chi won't route OPTIONS to GET; status isn't important here.

		mfs := scrapeMetrics(t, r)
		mf := mfs["webrpc_requests_total"]
		// We skipped this request entirely, so there should be no series with origin/no-header.example.
		assert.False(t, metricHasLabels(mf, map[string]string{"origin": "no-header.example"}))
	})
}

func scrapeMetrics(t *testing.T, r http.Handler) map[string]*dto.MetricFamily {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var p expfmt.TextParser
	mfs, err := p.TextToMetricFamilies(bytes.NewReader(rr.Body.Bytes()))
	assert.NoError(t, err)

	return mfs
}

func metricHasLabels(mf *dto.MetricFamily, labels map[string]string) bool {
	if mf == nil {
		return false
	}
	for _, m := range mf.GetMetric() {
		ok := true
		for wantName, wantValue := range labels {
			found := false
			for _, lp := range m.GetLabel() {
				if lp.GetName() == wantName && lp.GetValue() == wantValue {
					found = true
					break
				}
			}
			if !found {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
