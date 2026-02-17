package alert

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"testing"
)

func TestLogHandler_TriggersOnErrorf(t *testing.T) {
	var called bool
	handler := LogHandler(slog.NewTextHandler(io.Discard, nil), func(ctx context.Context, record slog.Record, err error) {
		called = true
	})
	logger := slog.New(handler)
	logger.Error("failed", slog.Any("error", Errorf("timeout")))
	if !called {
		t.Error("expected alert callback for Errorf")
	}
}

func TestLogHandler_UpgradesToLevelAlert(t *testing.T) {
	var logBuf strings.Builder
	baseHandler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{
		ReplaceAttr: ReplaceAttr(
			func(groups []string, a slog.Attr) slog.Attr { return a },
			slog.String("severity", "ALERT"),
		),
	})
	handler := LogHandler(baseHandler, func(ctx context.Context, record slog.Record, err error) {})
	logger := slog.New(handler)
	logger.Error("failed", slog.Any("error", Errorf("timeout")))
	logOut := logBuf.String()
	if !strings.Contains(logOut, "ALERT") {
		t.Errorf("expected log to contain severity=ALERT, got %q", logOut)
	}
}

func TestErrorf_CapturesStackFrames(t *testing.T) {
	err := Errorf("test error")
	var ae *alertError
	if !errors.As(err, &ae) {
		t.Fatal("expected *alertError from Errorf")
	}
	frames := ae.StackFrames()
	if len(frames) == 0 {
		t.Fatal("expected at least one stack frame")
	}
	var hasValidPC bool
	for _, pc := range frames {
		if pc != 0 {
			hasValidPC = true
			break
		}
	}
	if !hasValidPC {
		t.Fatal("expected at least one non-zero program counter in frames")
	}
	callers := runtime.CallersFrames(frames)
	frame, more := callers.Next()
	if !more && frame.PC == 0 {
		t.Fatal("could not resolve any stack frame")
	}
	if frame.Function == "" && frame.File == "" {
		t.Error("expected resolved frame to have Function or File")
	}
	if !strings.Contains(frame.Function, "Errorf") && !strings.Contains(frame.Function, "TestErrorf") {
		t.Errorf("expected frame to be Errorf or test caller, got %q", frame.Function)
	}
}

func TestLogHandler_DoesNotTriggerOnRegularError(t *testing.T) {
	var called bool
	handler := LogHandler(slog.NewTextHandler(io.Discard, nil), func(ctx context.Context, record slog.Record, err error) {
		called = true
	})
	logger := slog.New(handler)
	logger.Error("failed", slog.Any("error", fmt.Errorf("plain error")))
	if called {
		t.Error("expected no alert callback for regular error")
	}
}

func TestLogHandler_PassesWithAttrsAndGroups(t *testing.T) {
	var gotAttrs []slog.Attr
	handler := LogHandler(slog.NewTextHandler(io.Discard, nil), func(ctx context.Context, record slog.Record, err error) {
		record.Attrs(func(a slog.Attr) bool {
			gotAttrs = append(gotAttrs, a)
			return true
		})
	})

	logger := slog.New(handler).With(slog.String("service", "rpc")).WithGroup("meta")
	logger.Info("failed", slog.String("method", "ping"), slog.Any("error", Errorf("timeout")))

	if len(gotAttrs) != 1 || gotAttrs[0].Key != "meta" {
		t.Fatalf("expected one grouped attr with key=meta, got %+v", gotAttrs)
	}

	metaAttrs := gotAttrs[0].Value.Group()
	if len(metaAttrs) < 3 {
		t.Fatalf("expected grouped attrs to include service, method, error; got %+v", metaAttrs)
	}

	var gotService, gotMethod, gotError bool
	for _, a := range metaAttrs {
		switch a.Key {
		case "service":
			gotService = a.Value.String() == "rpc"
		case "method":
			gotMethod = a.Value.String() == "ping"
		case "error":
			e, ok := a.Value.Any().(error)
			if ok && e != nil {
				var ae *alertError
				gotError = errors.As(e, &ae)
			}
		}
	}
	if !gotService || !gotMethod || !gotError {
		t.Fatalf("expected grouped attrs with service=rpc method=ping and alert error, got %+v", metaAttrs)
	}
}

func TestLogHandler_PreventsRecursiveCallback(t *testing.T) {
	var callbackCalls int
	var out bytes.Buffer

	var logger *slog.Logger
	handler := LogHandler(slog.NewTextHandler(&out, nil), func(ctx context.Context, record slog.Record, err error) {
		callbackCalls++
		// Re-log the same alert error from inside callback.
		logger.ErrorContext(ctx, "callback relog", slog.Any("error", err))
	})
	logger = slog.New(handler)

	logger.ErrorContext(context.Background(), "root alert", slog.Any("error", Errorf("timeout")))

	if callbackCalls != 1 {
		t.Fatalf("expected callback to be called once, got %d", callbackCalls)
	}
}

func TestLogHandler_PanicsOnNilCallback(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when alertFn is nil")
		}
	}()
	_ = LogHandler(slog.NewTextHandler(io.Discard, nil), nil)
}
