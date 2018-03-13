package endpoint

import (
	"context"
	"io"
)

// Destination is the interface for destination endpoint.
type Destination interface {
	Fetchable() bool
	Writable() bool

	Write(ctx context.Context, path string, r io.ReadCloser) (err error)
	Fetch(ctx context.Context, path, url string) (err error)
	Dir(ctx context.Context, path string) (err error)
}
