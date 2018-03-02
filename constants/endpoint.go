package constants

// Endpoint the const for all supported endpoint.
const (
	EndpointQingStor = "qingstor"
	EndpointFs       = "fs"
	EndpointAliyun   = "aliyun"
)

// Constants for endpoint type.
const (
	SourceEndpoint uint8 = iota
	DestinationEndpoint
)
