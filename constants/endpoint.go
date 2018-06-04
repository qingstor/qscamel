package constants

// Endpoint the const for all supported endpoint.
const (
	EndpointAliyun   = "aliyun"
	EndpointFileList = "filelist"
	EndpointFs       = "fs"
	EndpointGCS      = "gcs"
	EndpointQingStor = "qingstor"
	EndpointQiniu    = "qiniu"
	EndpointS3       = "s3"
	EndpointUpyun    = "upyun"
)

// Constants for endpoint type.
const (
	SourceEndpoint uint8 = iota
	DestinationEndpoint
)
