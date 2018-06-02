package qingstor

// DirectoryContentType is the content type for qingstor directory.
const DirectoryContentType = "application/x-directory"

// MaxListObjectsLimit is the max limit for list objects.
const MaxListObjectsLimit = 1000

// Constants for storage class.
const (
	StorageClassStandard   = "STANDARD"
	StorageClassStandardIA = "STANDARD_IA"
)

// Multipart related constants.
// ref: https://docs.qingcloud.com/qingstor/api/object/multipart/index.html
const (
	// DefaultMultipartBoundarySize is the default multipart size.
	// 64 * 1024 * 1024 = 67108864 B = 64 MB
	DefaultMultipartSize = 67108864
	// DefaultMultipartBoundarySize is the default multipart boundary size.
	// 2 * 1024 * 1024 * 1024 = 2147483648 B = 2 GB
	DefaultMultipartBoundarySize = 2147483648
	// MaxMultipartBoundarySize is the max multipart boundary size.
	// Over this, put object will be reset by server.
	// 5 * 1024 * 1024 * 1024 = 5368709120 B = 5 GB
	MaxMultipartBoundarySize = 5368709120
)
