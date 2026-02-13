package middleware

import (
	"cmp"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	hostname, _ = os.Hostname()
	startedAt   = time.Now().UTC().Format(time.RFC3339)
)

// Health middleware responds with static JSON payload with the given version of the app
func Health(endpoint string, app string, version string) func(http.Handler) http.Handler {
	info := struct {
		App       string `json:"app"`
		Version   string `json:"version"`
		StartedAt string `json:"startedAt"`
		Hostname  string `json:"hostname"`
	}{
		App:       cmp.Or(app, os.Args[0], "unknown"),
		Version:   cmp.Or(version, "dev"),
		StartedAt: startedAt,
		Hostname:  cmp.Or(hostname, "unknown"),
	}
	resp, _ := json.Marshal(info)

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == endpoint {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
