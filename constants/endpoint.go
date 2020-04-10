package constants

// Endpoint the const for all supported endpoint.
const (
	EndpointAliyun   = "aliyun"
	EndpointAzblob   = "azblob"
	EndpointFileList = "filelist"
	EndpointFs       = "fs"
	EndpointGCS      = "gcs"
	EndpointHDFS     = "hdfs"
	EndpointQingStor = "qingstor"
	EndpointQiniu    = "qiniu"
	EndpointS3       = "s3"
	EndpointUpyun    = "upyun"
	EndpointCOS      = "cos"
)

// Constants for endpoint type.
const (
	SourceEndpoint uint8 = iota
	DestinationEndpoint
)
