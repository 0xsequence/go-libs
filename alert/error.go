package alert

import (
	"fmt"
	"runtime"
)

// Errorf creates a new error with a stack trace that triggers
// an alert when logged via alert.LogHandler.
func Errorf(format string, args ...any) error {
	err := &alertError{err: fmt.Errorf(format, args...)}
	runtime.Callers(1, err.frame.frames[:])
	return err
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
