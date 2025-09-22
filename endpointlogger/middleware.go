package endpointlogger

import (
	"net/http"
	"strings"
	"time"
)

// Middleware parse webrcp "service" "name" and endpoint name and adds it to context,
// also adds "since" so logs will have attribute when the log happened since received request to server
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		i := strings.LastIndex(path, "/")
		var endpoint, service string
		if i > 0 {
			endpoint = path[i+1:]

			// look for segment before last
			j := strings.LastIndex(path[:i], "/")
			if j == -1 {
				service = path[:i]
			} else {
				service = path[j+1 : i]
			}
		}

		if endpoint != "" && service != "" {
			r = r.WithContext(setValues(r.Context(), service, endpoint, time.Now().UTC()))
		}

		next.ServeHTTP(w, r)
	})
}
