package constants

import "errors"

// These errors can be returned when handle task.
var (
	// ErrTaskInvalid is returned when this task is invalid.
	ErrTaskInvalid = errors.New("task is invalid")
	// ErrTaskMismatch is returned when task content has been changed.
	ErrTaskMismatch = errors.New("task content mismatch")
	// ErrTaskNotFinished is returned when task is not finished.
	ErrTaskNotFinished = errors.New("task not finished")
	// ErrTaskNotFound is returned when task is not found.
	ErrTaskNotFound = errors.New("task not found")

	// ErrEndpointInvalid is returned when this endpoint is invalid.
	ErrEndpointInvalid = errors.New("endpoint is invalid")
	// ErrEndpointNotSupported is returned when this endpoint is not supported.
	ErrEndpointNotSupported = errors.New("endpoint type is not supported")
	// ErrEndpointFuncNotImplemented is return when a not implement function is called.
	ErrEndpointFuncNotImplemented = errors.New("endpoint does not implement this function")

	// ErrObjectTooLarge is returned when the object is too large.
	ErrObjectTooLarge = errors.New("object is too large")
	// ErrObjectInvalid is returned when the object is invalid.
	ErrObjectInvalid = errors.New("object is invalid")
)
