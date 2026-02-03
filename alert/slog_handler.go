package alert

import (
	"context"
	"errors"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

type Handler struct {
	next slog.Handler
}

// NewHandler wraps an existing slog.Handler and reports alert errors to Sentry.
func NewHandler(next slog.Handler) slog.Handler {
	return &Handler{next: next}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	// Always forward to the wrapped handler first
	err := h.next.Handle(ctx, r)

	var loggedErr error

	// Extract slog.ErrorKey ("error") if present
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == slog.ErrorKey {
			if e, ok := a.Value.Any().(error); ok {
				loggedErr = e
				return false
			}
		}
		return true
	})

	// Alert-worthy?
	if loggedErr != nil && errors.Is(loggedErr, ErrAlert) {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)

			// Copy slog attributes into Sentry extras
			r.Attrs(func(a slog.Attr) bool {
				scope.SetExtra(a.Key, a.Value.Any())
				return true
			})

			sentry.CaptureException(loggedErr)
		})
	}

	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{next: h.next.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{next: h.next.WithGroup(name)}
}
