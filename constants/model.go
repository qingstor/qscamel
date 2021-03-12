package constants

import (
	"fmt"
	"github.com/pengsrc/go-shared/buffer"
)

// Constants for task type.
const (
	TaskTypeCopy   = "copy"
	TaskTypeDelete = "delete"
	TaskTypeFetch  = "fetch"
)

// Constants for task status.
const (
	TaskStatusCreated  = "created"
	TaskStatusRunning  = "running"
	TaskStatusFinished = "finished"
)

// Constants for task ignore existing config.
const (
	TaskIgnoreExistingLastModified = "last_modified"
	TaskIgnoreExistingMD5Sum       = "md5sum"
)

// Constants for object types.
const (
	ObjectTypeDirectory = "directory"
	ObjectTypeSingle    = "single"
	ObjectTypePartial   = "partial"
)

// Constants for database key.
const (
	KeyTaskPrefix = "t:"

	// ObjectPrefixKey `~` is bigger than all ascii printable characters.
	ObjectPrefixKey = "~"

	KeyDirectoryObjectPrefix = "do:"
	KeySingleObjectPrefix    = "so:"
	KeyPartialObjectPrefix   = "po:"
)

// FormatTaskKey will format a task key.
func FormatTaskKey(t string) []byte {
	return []byte(KeyTaskPrefix + t)
}

// FormatDirectoryObjectKey will format a directory object key.
func FormatDirectoryObjectKey(t, s string) []byte {
	buf := buffer.GlobalBytesPool().Get()
	defer buf.Free()

	buf.AppendString(ObjectPrefixKey)
	buf.AppendString(t)
	buf.AppendString(":")
	buf.AppendString(KeyDirectoryObjectPrefix)
	buf.AppendString(s)

	b := make([]byte, buf.Len())
	copy(b, buf.Bytes())
	return b
}

// FormatSingleObjectKey will format a single object key.
func FormatSingleObjectKey(t, s string) []byte {
	buf := buffer.GlobalBytesPool().Get()
	defer buf.Free()

	buf.AppendString(ObjectPrefixKey)
	buf.AppendString(t)
	buf.AppendString(":")
	buf.AppendString(KeySingleObjectPrefix)
	buf.AppendString(s)

	b := make([]byte, buf.Len())
	copy(b, buf.Bytes())
	return b
}

// FormatPartialObjectKey will format a partial object key.
func FormatPartialObjectKey(t, s string, partNumber int) []byte {
	buf := buffer.GlobalBytesPool().Get()
	defer buf.Free()

	buf.AppendString(ObjectPrefixKey)
	buf.AppendString(t)
	buf.AppendString(":")
	buf.AppendString(KeyPartialObjectPrefix)
	// If s is empty, this key will be prefix for all partial object.
	if s != "" {
		buf.AppendString(s)
		buf.AppendString(":")
	}
	// If part number is lower than 0, this key will be prefix for
	// partial object with key.
	if partNumber >= 0 {
		// PartNumber must be [0,10000], so we can prefix enough "0"
		// to make sure the bytes order is correct.
		buf.AppendString(fmt.Sprintf("%05d", partNumber))
	}

	b := make([]byte, buf.Len())
	copy(b, buf.Bytes())
	return b
}
