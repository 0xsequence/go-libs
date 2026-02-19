# alert

Package `alert` adds alert semantics on top of `log/slog`:

- mark important errors with `alert.Errorf(...)`
- intercept those errors in a handler wrapper (`alert.LogHandler`)
- run a callback (for paging, Sentry, metrics, etc.)
- upgrade the log level to `alert.LevelAlert`
- map that level to sink-specific attrs with `alert.ReplaceAttr`

This keeps call-sites idiomatic `slog` while still producing GCP-native severity.
The package defines alert behavior; output-schema mapping is applied when wiring handlers.

## Why this package exists

`slog` has native levels and attributes. GCP Cloud Logging expects specific values in the
`severity` field (including `ALERT`).

`alert.LevelAlert` gives you a distinct semantic level in `slog`.
At logger wiring time, `alert.ReplaceAttr` can be composed into handler options to
translate that level into a sink-compatible field/value.

## Example usage

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/0xsequence/go-libs/alert"
	"github.com/go-chi/httplog/v3"
)

func main() {
	concise := true
	gcpSchema := httplog.SchemaGCP.Concise(concise).ReplaceAttr

	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   true,
		ReplaceAttr: alert.ReplaceAttr(gcpSchema, slog.String("severity", "ALERT")),
	})

	handler := alert.LogHandler(baseHandler, func(ctx context.Context, record slog.Record, err error) {
		// Send error to Sentry/PagerDuty/webhook/metrics etc.
		// You can collect record.Attrs() and record.Source().
		// IMPORTANT: do NOT log with slog here (use external side effects only).
	})

	logger := slog.New(handler)

	// Regular error: no alert callback, regular error severity.
	logger.Error("request failed", slog.Any("error", os.ErrNotExist))

	// Alert error: callback triggered + record level upgraded to LevelAlert,
	// then ReplaceAttr emits severity="ALERT" for GCP.
	logger.Error("critical timeout", slog.Any("error", alert.Errorf("timeout contacting upstream")))
}
```

Examples:

- GCP: `alert.ReplaceAttr(baseReplaceAttr, slog.String("severity", "ALERT"))`
- OTel-style: `alert.ReplaceAttr(baseReplaceAttr, slog.String("severityText", "ALERT"))`
- ECS-style: `alert.ReplaceAttr(baseReplaceAttr, slog.String("log.level", "alert"))`

## Sentry example (forward `error` + attrs)

This is possible with the current `LogHandler` callback. The callback receives the matched
`error` and a `slog.Record` whose attrs include call-site attrs and `logger.With(...)` context.

```go
import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0xsequence/go-libs/alert"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/traceid"
)

func sentryAlertHandler(base slog.Handler) slog.Handler {
	return alert.LogHandler(base, func(ctx context.Context, record slog.Record, err error) {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetContext(record.Message, map[string]any{
				"level":   "alert",
				"traceId": traceid.FromContext(ctx),
			})

			if source := record.Source(); source != nil {
				scope.SetTag("file_line", fmt.Sprintf("%s:%d", source.File, source.Line))
				scope.SetTag("func", source.Function)
			}

			record.Attrs(func(a slog.Attr) bool {
				scope.SetTag(a.Key, a.Value.String())
				return true
			})

			sentry.CaptureException(err)
		})
	})
}
```

## Error helpers

- `alert.Errorf(format, args...)`: create a new alert error with a formatted message.
- `alert.Error(err)`: wrap an existing error as an alert error.
- `alert.ErrorSkip(skip, err)`: advanced helper for wrapper packages that need caller-accurate stack frames (for example, `xlog.Alert`-style helpers).

## Operational notes

- Only errors created with `alert.Errorf(...)` trigger the alert callback.
- The callback receives `record` + `err`; `record.Attrs(...)` includes call-site and `logger.With(...)` attrs.
- `LevelAlert` is higher than `slog.LevelError`, so existing level filters still pass it.
- Callback rule: treat `alertFn` as a side-effect hook (Sentry, paging, webhooks, metrics).
- Do NOT log with slog from inside `alertFn` (especially alert errors), or you can trigger recursion.
