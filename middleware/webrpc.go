package middleware

import (
	"net/http"

	"github.com/0xsequence/go-libs/middleware/webrpc"
)

// Deprecated: Use webrpc.Telemetry(webrpc.Opts{Origin: false}) instead.
func WebrpcTelemetry(next http.Handler) http.Handler {
	return webrpc.Telemetry(webrpc.Opts{Origin: false})(next)
}
