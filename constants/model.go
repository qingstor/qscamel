package constants

// Constants for task type.
const (
	TaskTypeCopy  = "copy"
	TaskTypeFetch = "fetch"
)

// Constants for task status.
const (
	TaskStatusCreated  = "created"
	TaskStatusRunning  = "running"
	TaskStatusFinished = "finished"
)

// Constants for task ignore existing config.
const (
	TaskIgnoreExistingDisable     = "disable"
	TaskIgnoreExistingSize        = "size"
	TaskIgnoreExistingQuickMD5Sum = "quick_md5sum"
	TaskIgnoreExistingFullMD5Sum  = "full_md5sum"
)

// Constants for database key.
const (
	// prefixKey ~ is bigger than all ascii printable characters.
	prefixKey = "~"

	KeyTaskPrefix = "t:"

	KeyJobPrefix = "j:"

	KeyObjectPrefix = "o:"
)

// FormatTaskKey will format a task key.
func FormatTaskKey(t string) []byte {
	return []byte(KeyTaskPrefix + t)
}

// FormatJobKey will format a job key.
func FormatJobKey(t, s string) []byte {
	return []byte(prefixKey + t + ":" + KeyJobPrefix + s)
}

// FormatObjectKey will format a object key.
func FormatObjectKey(t, s string) []byte {
	return []byte(prefixKey + t + ":" + KeyObjectPrefix + s)
}
