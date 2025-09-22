package endpointlogger

import (
	"context"
	"log/slog"
	"time"
)

// LogHandler creates a new slog handler that will add "service", "endpoint" and "since" attributes to logs
func LogHandler(handler slog.Handler) slog.Handler {
	return &endpointHandler{
		handler: handler,
	}
}

type endpointHandler struct {
	handler slog.Handler
}

func (h *endpointHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *endpointHandler) Handle(ctx context.Context, record slog.Record) error {
	values, ok := getValues(ctx)
	if ok {
		record.AddAttrs(
			slog.String("service", values.service),
			slog.String("endpoint", values.endpoint),
			slog.Duration("since", time.Since(values.time)),
		)
	}

	return h.handler.Handle(ctx, record) //nolint
}

func (h *endpointHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &endpointHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *endpointHandler) WithGroup(name string) slog.Handler {
	return &endpointHandler{handler: h.handler.WithGroup(name)}
}
