package alert

import "fmt"

// ErrAlert is a sentinel error used only as a marker.
// It must never be returned directly.
var ErrAlert = &alertError{}

type alertError struct{}

func (e *alertError) Error() string {
	return "alert"
}

// Errorf marks an error as alert-worthy while preserving %w semantics.
//
// Usage:
//   return alert.Errorf("failed to load user %d: %w", id, err)
func Errorf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrAlert}, args...)...)
}
