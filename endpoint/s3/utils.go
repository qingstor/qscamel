package s3

import (
	"github.com/yunify/qscamel/constants"
)

// calculatePartSize will calculate the object's part size.
func calculatePartSize(size int64) (partSize int64, err error) {
	partSize = DefaultMultipartSize

	for size/partSize >= int64(MaxMultipartNumber) {
		if partSize < MaxAutoMultipartSize {
			partSize = partSize << 1
			continue
		}
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		partSize = size/int64(MaxMultipartNumber) + 1
		break
	}

	if partSize > MaxMultipartBoundarySize {
		err = constants.ErrObjectTooLarge
		return
	}

	return
}
