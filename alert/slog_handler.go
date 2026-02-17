package alert

import (
	"context"
	"errors"
	"log/slog"
)

// LevelAlert is a custom slog level (16) for alert errors. It is greater than
// LevelError (8), so handlers with Level >= LevelError will pass it through.
const LevelAlert = slog.Level(16)

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
// The callback receives the context, the full record (including record.Source()
// for the slog caller's file/line when AddSource is enabled), and the matched error.
//
// Example (simple, no dependencies):
//
//	slogHandler = alert.LogHandler(slogHandler, func(ctx context.Context, record slog.Record, err error) {
//	    if source := record.Source(); source != nil {
//	        fmt.Printf("ALERT %s:%d: %v\n", source.File, source.Line, err)
//	    } else {
//	        fmt.Printf("ALERT: %v\n", err)
//	    }
//	})
//
// Example (Sentry, capturing the slog caller as the exception location):
//
//	slogHandler = alert.LogHandler(slogHandler, func(ctx context.Context, record slog.Record, err error) {
//	    sentry.WithScope(func(scope *sentry.Scope) {
//	        if source := record.Source(); source != nil {
//	            scope.SetTag("log_caller", fmt.Sprintf("%s:%d", source.File, source.Line))
//	            scope.SetExtra("log_source", map[string]any{
//	                "file":     source.File,
//	                "line":     source.Line,
//	                "function": source.Function,
//	            })
//	        }
//	        sentry.CaptureException(err)
//	    })
//	})
func LogHandler(handler slog.Handler, alertFn func(ctx context.Context, record slog.Record, err error)) slog.Handler {
	return &alertHandler{
		handler: handler,
		alertFn: alertFn,
	}
}

type alertHandler struct {
	handler slog.Handler
	alertFn func(ctx context.Context, record slog.Record, err error)
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
		h.alertFn(ctx, record, alertErr)
		return h.handler.Handle(ctx, record) //nolint:wrapcheck
	}
	return h.handler.Handle(ctx, record) //nolint:wrapcheck
}

func (h *alertHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &alertHandler{handler: h.handler.WithAttrs(attrs), alertFn: h.alertFn}
}

func (h *alertHandler) WithGroup(name string) slog.Handler {
	return &alertHandler{handler: h.handler.WithGroup(name), alertFn: h.alertFn}
}
