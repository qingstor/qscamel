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

	List(ctx context.Context, j *model.Job, rc chan *model.Object)
	Read(ctx context.Context, p string) (r io.ReadCloser, err error)
	Reach(ctx context.Context, p string) (url string, err error)
}
