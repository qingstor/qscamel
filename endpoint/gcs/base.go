package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "gcs:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string, _ bool) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	o := c.client.Object(cp)

	return o.NewReader(ctx)
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	r, err = c.client.Object(cp).NewRangeReader(ctx, offset, size)
	if err != nil {
		return
	}

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string, _ bool) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.Object(cp).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, nil
		}
		return
	}
	o = &model.SingleObject{
		Key:          p,
		Size:         resp.Size,
		LastModified: resp.Updated.Unix(),
		MD5:          string(resp.MD5),
	}
	return
}
