package azblob

import (
	"context"
	"errors"
	"io"

	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types/pairs"

	"github.com/yunify/qscamel/model"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "azblob:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.Reader, err error) {
	return c.client.ReadWithContext(ctx, p)
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	return c.client.ReadWithContext(ctx, p, pairs.WithOffset(offset), pairs.WithSize(size))
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.SingleObject, err error) {
	so, err := c.client.StatWithContext(ctx, p)
	if err != nil {
		if errors.Is(err, services.ErrObjectNotExist) {
			return nil, nil
		}
		return
	}

	o = &model.SingleObject{
		Key:          p,
		Size:         so.Size,
		LastModified: so.UpdatedAt.Unix(),
	}

	if v, ok := so.GetContentMD5(); ok {
		o.MD5 = v
	}
	return
}
