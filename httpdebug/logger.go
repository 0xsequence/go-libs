package httpdebug

import (
	"context"
	"log/slog"
)

// LogHandler creates a new slog handler that will emit log records
// at all levels, if debug mode is enabled in the context.
func LogHandler(debugHeader Header) func(handler slog.Handler) slog.Handler {
	return func(handler slog.Handler) slog.Handler {
		return &debugHandler{
			next: handler,
			h:    debugHeader,
		}
	}
}

type debugHandler struct {
	next slog.Handler
	h    Header
}

func (h *debugHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if IsDebugModeEnabled(ctx) {
		return true
	}

	return h.next.Enabled(ctx, level)
}

func (h *debugHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.next.Handle(ctx, record) //nolint
}

func (h *debugHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &debugHandler{next: h.next.WithAttrs(attrs), h: h.h}
}

func (h *debugHandler) WithGroup(name string) slog.Handler {
	return &debugHandler{next: h.next.WithGroup(name), h: h.h}
}
