package middleware

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/metrics"

	"github.com/0xsequence/marketplace-api/lib/status/proto"
)

var (
	// Total number of requests versioned with Webrpc header.
	requestsTotal = metrics.CounterWith[labels]("webrpc_requests_total", "Total number of webrpc client requests.")
)

type labels struct {
	Client string `label:"client"`
	Status string `label:"status"`
}

// WebrpcTelemetry is a middleware that extracts webrpc client information from request headers,
// logs it to request log for traceability, and collects usage metrics for API analytics.
func WebrpcTelemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version, err := proto.VersionFromHeader(r.Header)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		webrpcClient := version.SchemaName + "@" + version.SchemaVersion

		httplog.SetAttrs(r.Context(), slog.String("webrpcClient", webrpcClient))

		ww, ok := w.(middleware.WrapResponseWriter)
		if !ok {
			ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		}

		defer func() {
			requestsTotal.Inc(labels{
				Client: webrpcClient,
				Status: strconv.Itoa(ww.Status()),
			})
		}()

		next.ServeHTTP(ww, r)
	})
}
