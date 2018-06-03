package constants

import "errors"

// These errors can be returned when handle task.
var (
	// ErrTaskMismatch is returned when task content has been changed.
	ErrTaskMismatch = errors.New("task content mismatch")
	// ErrTaskNotFinished is returned when task is not finished.
	ErrTaskNotFinished = errors.New("task not finished")
	// ErrEndpointInvalid is returned when this endpoint is invalid.
	ErrEndpointInvalid = errors.New("endpoint is invalid")
	// ErrEndpointNotSupported is returned when this endpoint is not supported.
	ErrEndpointNotSupported = errors.New("endpoint type is not supported")
	// ErrEndpointFuncNotImplemented is return when a not implement function is called.
	ErrEndpointFuncNotImplemented = errors.New("endpoint does not implement this function")
)
