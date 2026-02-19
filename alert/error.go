package alert

import (
	"fmt"
	"runtime"
)

// Errorf creates a new error with a stack trace that triggers
// an alert when logged via alert.LogHandler.
func Errorf(format string, args ...any) error {
	return newAlertError(1, fmt.Errorf(format, args...))
}

// Error wraps an existing error with alert semantics and captures
// a stack trace at the call site.
func Error(err error) error {
	return newAlertError(1, err)
}

// ErrorSkip wraps an existing error with alert semantics and captures
// stack frames while skipping additional caller frames.
// Use this when creating helper wrappers in another package.
func ErrorSkip(skip int, err error) error {
	return newAlertError(1+skip, err)
}

func newAlertError(skip int, err error) error {
	alertErr := &alertError{err: err}
	runtime.Callers(1+skip, alertErr.frame.frames[:])
	return alertErr
}

// alertError triggers alerts when logged.
type alertError struct {
	err   error
	frame struct{ frames [3]uintptr }
}

func (e *alertError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "alert"
}

func (e *alertError) Unwrap() error {
	return e.err
}

func (e alertError) StackFrames() []uintptr {
	return e.frame.frames[:]
}
