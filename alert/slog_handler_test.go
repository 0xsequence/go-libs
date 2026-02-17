package alert

import (
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
