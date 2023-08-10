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
	Stat(ctx context.Context, p string, isDir bool) (o *model.SingleObject, err error)

	// Read will return a reader.
	Read(ctx context.Context, p string, isDir bool) (r io.Reader, err error)
	// ReadRange will read content with range [offset, offset+size)
	ReadRange(ctx context.Context, p string, offset, size int64) (r io.Reader, err error)
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

	// InitPart will inti a multipart upload.
	InitPart(ctx context.Context, p string, size int64, meta map[string]string) (uploadID string, partSize int64, partNumbers int, err error)
	// UploadPart will upload a part.
	UploadPart(ctx context.Context, o *model.PartialObject, r io.Reader) (err error)
	CompleteParts(ctx context.Context, path string, uploadId string, totalNumber int) (err error)
	// Partable will return whether current endpoint supports multipart upload.
	Partable() bool

	// Write will read data from the reader and write to endpoint.
	Write(ctx context.Context, path string, size int64, r io.Reader, isDir bool, meta map[string]string) (err error)
	// Writable will return whether current endpoint supports write.
	Writable() bool
}

// Source is the interface for source endpoint.
type Source interface {
	Base

	// List will list from the job.
	List(ctx context.Context, j *model.DirectoryObject, fn func(model.Object)) (err error)

	// Reach will return an accessible url.
	Reach(ctx context.Context, p string) (url string, err error)
	// Reachable will return whether current endpoint supports reach.
	Reachable() bool
}
