package aliyun

import (
	"context"
	"io"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "aliyun:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	r, err = c.client.GetObject(cp)
	if err != nil {
		return
	}

	return
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	r, err = c.client.GetObject(cp, oss.Range(offset, offset+size-1))
	if err != nil {
		return
	}

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObjectMeta(cp)
	if err != nil {
		if e, ok := err.(*oss.ServiceError); ok {
			// If object not found, we just need to return a nil object.
			if e.StatusCode == 404 {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}

	// Parse content length.
	size, err := strconv.ParseInt(resp.Get("Content-Length"), 10, 64)
	if err != nil {
		logrus.Errorf("Content length parsed failed for %v.", err)
		return
	}
	// Parse last modified.
	lastModified := convert.StringToTimestamp(resp.Get("Last-Modified"), convert.RFC822)

	o = &model.SingleObject{
		Key:          p,
		Size:         size,
		LastModified: lastModified,
		MD5:          resp.Get("ETag"),
	}
	return
}
