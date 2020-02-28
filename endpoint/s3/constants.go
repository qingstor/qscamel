package s3

// MaxKeys is the max limit for list objects.
const MaxKeys = 1000

// ErrorCodeNotFound is the error code for key not found.
const ErrorCodeNotFound = "NoSuchKey"

// MaxListObjectsLimit is the max limit for list objects.
const MaxListObjectsLimit = 1000

// Multipart related constants.
// ref: https://docs.qingcloud.com/qingstor/api/object/multipart/index.html
const (
	// DefaultMultipartBoundarySize is the default multipart size.
	// 64 * 1024 * 1024 = 67108864 B = 64 MB
	DefaultMultipartSize = 67108864
	// MaxAutoMultipartSize is the max auto multipart size.
	// If part size is over MaxAutoMultipartSize, we will not detect it any more.
	// 1024 * 1024 * 1024 = 1073741824 B = 1 GB
	MaxAutoMultipartSize = 1073741824
	// MaxMultipartNumber is the max part that QingStor supported.
	MaxMultipartNumber = 10000
	// MaxMultipartBoundarySize is the max multipart boundary size.
	// Over this, put object will be reset by server.
	// 5 * 1024 * 1024 * 1024 = 5368709120 B = 5 GB
	MaxMultipartBoundarySize = 5368709120
)
