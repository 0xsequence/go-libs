package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/metrics"
)

var (
	// Total number of requests versioned with Webrpc header.
	requestsTotal = metrics.CounterWith[labels]("webrpc_requests_total", "Total number of webrpc client requests.")
)

type labels struct {
	Gen    string `label:"gen"`
	Schema string `label:"schema"`
	Status string `label:"status"`
}

// WebrpcTelemetry is a middleware that extracts webrpc client information from request headers,
// logs it to request log for traceability, and collects usage metrics for API analytics.
func WebrpcTelemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webrpcHeader := r.Header.Get("Webrpc")
		if webrpcHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Webrpc: webrpc@v0.25.1;gen-golang@v0.19.0;marketplace-api@v25.9.1;...
		versions := strings.Split(webrpcHeader, ";")
		if len(versions) < 3 {
			next.ServeHTTP(w, r)
			return
		}

		webrpcGen, _, _ := strings.Cut(versions[1], "@") // gen-golang@v0.19.0 -> gen-golang
		webrpcSchema := versions[2]                      // marketplace-api@v25.9.1

		httplog.SetAttrs(r.Context(),
			slog.String("webrpcGen", webrpcGen),
			slog.String("webrpcSchema", webrpcSchema),
		)

		ww, ok := w.(middleware.WrapResponseWriter)
		if !ok {
			ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		}

		defer func() {
			requestsTotal.Inc(labels{
				Gen:    webrpcGen,
				Schema: webrpcSchema,
				Status: strconv.Itoa(ww.Status()),
			})
		}()

		next.ServeHTTP(ww, r)
	})
}
