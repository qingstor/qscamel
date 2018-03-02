package endpoint

import (
	"context"
	"io"
)

// Destination is the interface for destination endpoint.
type Destination interface {
	Fetchable() bool
	Writable() bool

	Write(ctx context.Context, path string, r io.Reader) (err error)
	Fetch(ctx context.Context, path string) (err error)
	Dir(ctx context.Context, path string) (err error)
}
