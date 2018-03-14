package endpoint

import (
	"context"
	"io"

	"github.com/yunify/qscamel/model"
)

// Source is the interface for source endpoint.
type Source interface {
	Reachable() bool
	Readable() bool

	List(ctx context.Context, p string, rc chan *model.Object) (err error)
	Read(ctx context.Context, p string) (r io.ReadCloser, err error)
	Reach(ctx context.Context, p string) (url string, err error)
}
