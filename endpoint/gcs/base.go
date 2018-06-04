package gcs

import (
	"context"
	"io"
	"strings"

	"cloud.google.com/go/storage"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "gcs:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	o := c.client.Object(cp)

	return o.NewReader(ctx)
}

// ReadAt implement source.ReadAt
func (c *Client) ReadAt(
	ctx context.Context, p string, start, end int64,
) (b []byte, err error) {
	cp := utils.Join(c.Path, p)

	r, err := c.client.Object(cp).NewRangeReader(ctx, start, end-start+1)
	if err != nil {
		return
	}
	defer r.Close()

	b = make([]byte, end-start+1)
	_, err = r.Read(b)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.Object(cp).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, nil
		}
		return
	}
	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         resp.Size,
		LastModified: resp.Updated.Unix(),
		MD5:          string(resp.MD5),
	}
	return
}
