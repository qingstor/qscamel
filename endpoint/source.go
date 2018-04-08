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
	MD5able() bool

	List(ctx context.Context, j *model.Job, fn func(*model.Object)) (err error)

	Read(ctx context.Context, p string) (r io.ReadCloser, err error)
	Reach(ctx context.Context, p string) (url string, err error)

	Stat(ctx context.Context, p string) (o *model.Object, err error)

	MD5(ctx context.Context, p string) (b string, err error)
}
