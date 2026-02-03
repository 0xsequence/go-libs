package webrpc

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/metrics"
)

// Total number of requests.
var requestsTotal = metrics.CounterWith[labels]("webrpc_requests_total", "Total number of webrpc requests.")

type labels struct {
	Gen    string `label:"gen"`
	Schema string `label:"schema"`
	Status string `label:"status"`
	Origin string `label:"origin"`
}

type Opts struct {
	// Track origin label in metrics.
	// NOTE: Cardinality grows with the number of unique origin headers.
	Origin bool

	Skip func(r *http.Request) bool
}

// Telemetry is a middleware that extracts webrpc client information from request headers,
// logs it to request log for traceability, and collects usage metrics for API analytics.
func Telemetry(opts Opts) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if opts.Skip != nil && opts.Skip(r) {
				next.ServeHTTP(w, r)
				return
			}

			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			var labels labels
			defer func() {
				labels.Status = strconv.Itoa(ww.Status())
				requestsTotal.Inc(labels)
			}()

			webrpcGen, webrpcSchema := parseWebrpcHeader(r.Header.Get("Webrpc"))
			if webrpcSchema != "" {
				labels.Gen = webrpcGen
				labels.Schema = webrpcSchema
				httplog.SetAttrs(r.Context(),
					slog.String("webrpcGen", webrpcGen),
					slog.String("webrpcSchema", webrpcSchema),
				)
			}

			if opts.Origin {
				if origin := strings.TrimSpace(r.Header.Get("Origin")); origin != "" && origin != "null" {
					if u, err := url.Parse(origin); err == nil && u.Scheme != "" && u.Host != "" {
						labels.Origin = u.Host
					}
				}
			}

			next.ServeHTTP(ww, r)
		})
	}
}

func parseWebrpcHeader(header string) (string, string) {
	versions := strings.Split(header, ";")
	if len(versions) < 3 {
		return "", ""
	}
	webrpcGen, _, _ := strings.Cut(versions[1], "@") // gen-golang@v0.19.0 -> gen-golang
	webrpcSchema := versions[2]                      // marketplace-api@v25.9.1
	return webrpcGen, webrpcSchema
}
