package endpoint

import (
	"context"
	"github.com/yunify/qscamel/model"
	"io"
)

// Source is the interface for source endpoint.
type Source interface {
	Reachable() bool
	Readable() bool

	List(ctx context.Context, p string) (o []model.Object, err error)
	Read(ctx context.Context, p string) (r io.ReadCloser, err error)
}
