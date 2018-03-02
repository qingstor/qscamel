package constants

import "errors"

// These errors can be returned when handle task.
var (
	// ErrTaskMismatch is returned when task content has been changed.
	ErrTaskMismatch = errors.New("task content mismatch")
)
