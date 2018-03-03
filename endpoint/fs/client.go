package fs

import (
	"context"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the struct for POSIX file system endpoint.
type Client struct {
	Path string
}

// New will create a Fs.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	c = &Client{}

	// Set prefix.
	c.Path = e.Path

	return
}
