package hdfs

import (
	"context"
	"io"
	"os"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "hdfs:" + c.Address
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	fr, err := c.client.Open(cp)
	if err != nil {
		return
	}
	return fr, nil
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	fr, err := c.client.Open(cp)
	if err != nil {
		return
	}

	r = io.NewSectionReader(fr, offset, size)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)

	fi, err := c.client.Stat(cp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return
	}

	// We will not calculate md5 while stating object.
	o = &model.SingleObject{
		Key:          p,
		Size:         fi.Size(),
		LastModified: fi.ModTime().Unix(),
	}
	return
}
