package endpointlogger

import (
	"context"
	"time"
)

type ctxKey struct{}

type ctxVal struct {
	service  string
	endpoint string
	time     time.Time
}

func setValues(ctx context.Context, service, endpoint string, t time.Time) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxVal{service: service, endpoint: endpoint, time: t})
}

func getValues(ctx context.Context) (ctxVal, bool) {
	if v, ok := ctx.Value(ctxKey{}).(ctxVal); ok {
		return v, ok
	}

	return ctxVal{}, false
}
