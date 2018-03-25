package constants

import "github.com/pengsrc/go-shared/buffer"

// Constants for task type.
const (
	TaskTypeCopy        = "copy"
	TaskTypeFetch       = "fetch"
	TaskTypeVerifyCopy  = "verify+copy"
	TaskTypeVerifyFetch = "verify+fetch"
)

// Constants for task status.
const (
	TaskStatusRunning  = "running"
	TaskStatusFinished = "finished"
)

// Constants for database key.
const (
	KeyTaskList   = "t"
	KeyTaskPrefix = "t:"

	KeyJobPrefix = "j:"

	KeyObjectPrefix = "o:"
)

// FormatTaskKey will format a task key.
func FormatTaskKey(s string) []byte {
	return []byte(KeyTaskPrefix + s)
}

// FormatJobKey will format a job key.
func FormatJobKey(u uint64) []byte {
	b := buffer.GlobalBytesPool().Get()
	defer b.Free()

	b.AppendString(KeyJobPrefix)
	b.AppendUint(u)

	return b.Bytes()
}

// FormatObjectKey will format a object key.
func FormatObjectKey(name string) []byte {
	b := buffer.GlobalBytesPool().Get()
	defer b.Free()

	b.AppendString(KeyObjectPrefix)
	b.AppendString(name)

	return b.Bytes()
}
