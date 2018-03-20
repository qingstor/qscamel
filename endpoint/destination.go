package endpoint

import (
	"context"
	"io"

	"github.com/yunify/qscamel/model"
)

// Destination is the interface for destination endpoint.
type Destination interface {
	Fetchable() bool
	Writable() bool

	Write(ctx context.Context, path string, r io.ReadCloser) (err error)
	Fetch(ctx context.Context, path, url string) (err error)

	Stat(ctx context.Context, p string) (o *model.Object, err error)
}
