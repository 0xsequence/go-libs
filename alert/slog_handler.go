package alert

import (
	"context"
	"errors"
	"log/slog"
)

// LevelAlert is a custom slog level (16) for alert errors. It is greater than
// LevelError (8), so handlers with Level >= LevelError will pass it through.
const LevelAlert = slog.Level(16)

type callbackContextKey struct{}

// ReplaceAttr wraps a ReplaceAttr function to map LevelAlert to a caller-
// provided attr. This keeps alert semantics in this package while letting callers
// pick sink-specific fields/values (for example, GCP severity).
//
//	replaceAttr := alert.ReplaceAttr(
//		httplog.SchemaGCP.Concise(concise).ReplaceAttr,
//		slog.String("severity", "ALERT"),
//	)
func ReplaceAttr(next func(groups []string, a slog.Attr) slog.Attr, alertAttr slog.Attr) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) == 0 && a.Key == slog.LevelKey {
			if level, ok := a.Value.Any().(slog.Level); ok && level == LevelAlert {
				return alertAttr
			}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

// LogHandler wraps a slog.Handler and invokes the alert callback when a log
// record contains an error from Errorf (private alertError type). When an alert
// is triggered, the record's level is upgraded to LevelAlert before passing to
// the next handler, so GCP receives severity="ALERT" and level filters treat it
// as >= ERROR. Use ReplaceAttr to map LevelAlert to sink-specific attrs.
//
// The callback receives the matched error and a record that includes call-site
// attrs plus logger.With(...) / logger.WithGroup(...) context.
//
// Example (simple, no dependencies):
//
//	slogHandler = alert.LogHandler(slogHandler, func(ctx context.Context, record slog.Record, err error) {
//	    record.Attrs(func(a slog.Attr) bool {
//	        fmt.Println(a.Key, a.Value.String())
//	        return true
//	    })
//	})
//
// Example (Sentry, capturing the slog caller as the exception location):
//
//	slogHandler = alert.LogHandler(slogHandler, func(ctx context.Context, record slog.Record, err error) {
//	    sentry.CaptureException(err)
//	})
func LogHandler(handler slog.Handler, alertFn func(ctx context.Context, record slog.Record, err error)) slog.Handler {
	return &alertHandler{
		handler: handler,
		alertFn: alertFn,
	}
}

type alertHandler struct {
	handler    slog.Handler
	alertFn    func(ctx context.Context, record slog.Record, err error)
	parent     *alertHandler
	localAttrs []slog.Attr
	localGroup string
}

func (h *alertHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *alertHandler) Handle(ctx context.Context, record slog.Record) error {
	var alertErr error
	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "error" || a.Value.Kind() != slog.KindAny {
			return true
		}
		e, ok := a.Value.Any().(error)
		if !ok || e == nil {
			return true
		}

		var ae *alertError
		if errors.As(e, &ae) {
			alertErr = e
			return false
		}
		return true
	})
	if alertErr != nil {
		record.Level = LevelAlert
		if inCallback, _ := ctx.Value(callbackContextKey{}).(bool); inCallback {
			// Prevent infinite recursion when alertFn logs alert errors again
			// through the same handler chain.
			return h.handler.Handle(ctx, record) //nolint:wrapcheck
		}
		callbackRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
		callbackRecord.AddAttrs(h.buildAttrs(record)...)
		callbackCtx := context.WithValue(ctx, callbackContextKey{}, true)
		h.alertFn(callbackCtx, callbackRecord, alertErr)
		return h.handler.Handle(ctx, record) //nolint:wrapcheck
	}
	return h.handler.Handle(ctx, record) //nolint:wrapcheck
}

func (h *alertHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	localAttrs := append([]slog.Attr(nil), attrs...)
	return &alertHandler{
		handler:    h.handler.WithAttrs(attrs),
		alertFn:    h.alertFn,
		parent:     h,
		localAttrs: localAttrs,
	}
}

func (h *alertHandler) WithGroup(name string) slog.Handler {
	return &alertHandler{
		handler:    h.handler.WithGroup(name),
		alertFn:    h.alertFn,
		parent:     h,
		localGroup: name,
	}
}

// buildAttrs flattens logger.With(...) attrs from the parent chain and
// appends the current record attrs. If WithGroup was used, the final attrs are
// wrapped into nested slog groups in the same order they were applied.
func (h *alertHandler) buildAttrs(record slog.Record) []slog.Attr {
	var chain []*alertHandler
	for cur := h; cur != nil; cur = cur.parent {
		chain = append(chain, cur)
	}

	var groups []string
	var attrs []slog.Attr
	for i := len(chain) - 1; i >= 0; i-- {
		node := chain[i]
		if node.localGroup != "" {
			groups = append(groups, node.localGroup)
		}
		attrs = append(attrs, node.localAttrs...)
	}
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	if len(groups) == 0 {
		return attrs
	}

	grouped := slog.Attr{
		Key:   groups[len(groups)-1],
		Value: slog.GroupValue(attrs...),
	}
	for i := len(groups) - 2; i >= 0; i-- {
		grouped = slog.Attr{
			Key:   groups[i],
			Value: slog.GroupValue(grouped),
		}
	}
	return []slog.Attr{grouped}
}
