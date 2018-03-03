package constants

import (
	"net"
	"net/http"
	"time"
)

// DefaultClient is the default HTTP client for qscamel.
var DefaultClient = &http.Client{
	// Set timeout to 0 to prevent file closed too early.
	Timeout: 0,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			// With or without a timeout, the operating system may impose
			// its own earlier timeout
			Timeout: 1 * time.Minute,
			// Do not keep alive for too long.
			KeepAlive: 30 * time.Second,
			// XXX: DualStack enables RFC 6555-compliant "Happy Eyeballs" dialing
			// when the network is "tcp" and the destination is a host name
			// with both IPv4 and IPv6 addresses. This allows a client to
			// tolerate networks where one address family is silently broken
			DualStack: false,
		}).DialContext,
		MaxIdleConns:          0,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second, //Default
		ExpectContinueTimeout: 2 * time.Second,
	},
}
