package endpointlogger

import (
	"net/http"
	"strings"
	"time"
)

// Middleware captures webrpc service and endpoint name to a context
// also adds "since" so logs will have attribute when the log happened since received request to server
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var endpoint, service string
		parts := strings.Split(path, "/")
		if len(parts) > 3 {
			// check if the path is in webrpc format, and if so, then extract the service and endpoint
			if parts[len(parts)-3] == "rpc" {
				service = parts[len(parts)-2]
				endpoint = parts[len(parts)-1]
			}

			if endpoint != "" && service != "" {
				r = r.WithContext(setValues(r.Context(), service, endpoint, time.Now().UTC()))
			}
		}

		next.ServeHTTP(w, r)
	})
}
