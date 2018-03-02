package fs

import (
	"context"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

var (
	t *model.Task
)

// Fs is the struct for POSIX file system endpoint.
type Fs struct {
	Path string
}

// New will create a Fs.
func New(ctx context.Context, et uint8) (f *Fs, err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	f = &Fs{}

	// Set prefix.
	f.Path = e.Path

	return
}
