package httpdebug

import (
	"context"
)

type ctxKey struct{}

func isDebugModeEnabled(ctx context.Context) bool {
	return ctx.Value(ctxKey{}) != nil
}

func enableDebugMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, struct{}{})
}
