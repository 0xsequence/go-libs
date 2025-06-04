package httpdebug

import (
	"context"
	"log/slog"
)

type debugHandler struct {
	next slog.Handler
	cfg  Config
}

// LogHandler creates a new handler that wraps another handler
// and logs all levels if the debug context value is present in context.
func (c *Client) LogHandler(handler slog.Handler) slog.Handler {
	return &debugHandler{
		next: handler,
		cfg:  c.cfg,
	}
}

// If the debugging header was parsed, then the enable would return true, so it will log also Debug level
func (h *debugHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if IsDebug(ctx) {
		return true
	}

	return h.next.Enabled(ctx, level)
}

func (h *debugHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.next.Handle(ctx, record) //nolint
}

func (h *debugHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &debugHandler{next: h.next.WithAttrs(attrs), cfg: h.cfg}
}

func (h *debugHandler) WithGroup(name string) slog.Handler {
	return &debugHandler{next: h.next.WithGroup(name), cfg: h.cfg}
}
