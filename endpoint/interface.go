package endpoint

import (
	"context"
	"io"

	"github.com/yunify/qscamel/model"
)

// Base is the interface that both Source and Destination should implement.
type Base interface {
	// Name will return the endpoint's name.
	Name(ctx context.Context) (name string)

	// Stat will get the metadata.
	Stat(ctx context.Context, p string) (o *model.Object, err error)

	// Read will return a reader.
	Read(ctx context.Context, p string) (r io.ReadCloser, err error)
	// ReadAt will read content with range [start, end]
	ReadAt(ctx context.Context, p string, start, end int64) (b []byte, err error)
}

// Destination is the interface for destination endpoint.
type Destination interface {
	Base

	// Delete will use endpoint to delete the path.
	Delete(ctx context.Context, p string) (err error)
	// Deletable will return whether current endpoint supports delete.
	Deletable() bool

	// Fetch will use endpoint to fetch the url.
	Fetch(ctx context.Context, path, url string) (err error)
	// Fetchable will return whether current endpoint supports fetch.
	Fetchable() bool

	// Write will read data from the reader and write to endpoint.
	Write(ctx context.Context, path string, size int64, r io.ReadCloser) (err error)
	// Writable will return whether current endpoint supports write.
	Writable() bool
}

// Source is the interface for source endpoint.
type Source interface {
	Base

	// List will list from the job.
	List(ctx context.Context, j *model.Job, fn func(*model.Object)) (err error)

	// Reach will return an accessible url.
	Reach(ctx context.Context, p string) (url string, err error)
	// Reachable will return whether current endpoint supports reach.
	Reachable() bool
}
