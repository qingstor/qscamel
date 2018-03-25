package constants

import "errors"

// These errors can be returned when handle task.
var (
	// ErrTaskMismatch is returned when task content has been changed.
	ErrTaskMismatch = errors.New("task content mismatch")
	// ErrNotFinishedObject is returned when there are not finished object.
	ErrNotFinishedObject = errors.New("object not finished")
	// ErrEndpointInvalid is returned when this endpoint is invalid.
	ErrEndpointInvalid = errors.New("endpoint is invalid")
)
