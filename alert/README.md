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
		// Send to pager/Sentry/etc. You also get record.Source() when AddSource is enabled.
		slog.Default().WarnContext(ctx, "alert callback triggered", slog.Any("error", err))
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
alert error (`err`) and the full `slog.Record`, so you can send the exception and attach
context from the record attrs.

```go
import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0xsequence/go-libs/alert"
	"github.com/getsentry/sentry-go"
)

func sentryAlertHandler(base slog.Handler) slog.Handler {
	return alert.LogHandler(base, func(ctx context.Context, record slog.Record, err error) {
		sentry.WithScope(func(scope *sentry.Scope) {
			// Primary exception payload (the matched "error" attr from alert.Errorf).
			scope.SetContext("alert", map[string]any{
				"message": record.Message,
				"level":   record.Level.String(),
			})

			// Add caller info when AddSource is enabled.
			if source := record.Source(); source != nil {
				scope.SetTag("log.file", source.File)
				scope.SetTag("log.func", source.Function)
				scope.SetTag("log.line", fmt.Sprintf("%d", source.Line))
			}

			// Attach remaining attrs as Sentry extras.
			record.Attrs(func(a slog.Attr) bool {
				if a.Key == "error" {
					return true // already captured as exception payload
				}
				scope.SetExtra(a.Key, a.Value.String())
				return true
			})

			sentry.CaptureException(err)
		})
	})
}
```

## Operational notes

- Only errors created with `alert.Errorf(...)` trigger the alert callback.
- `LogHandler` preserves record attributes and source information.
- `LevelAlert` is higher than `slog.LevelError`, so existing level filters still pass it.
