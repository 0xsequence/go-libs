package httpdebug

import (
	"context"
	"net/http"
)

type ctxKey struct{}

func IsDebugModeEnabled(ctx context.Context) bool {
	return ctx.Value(ctxKey{}) != nil
}

func IsDebugHeaderSet(r *http.Request) bool {
	return IsDebugModeEnabled(r.Context())
}

func enableDebugMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, struct{}{})
}
