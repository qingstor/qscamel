package cos

import (
	"context"
	"fmt"
	"io"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "tencent cos:" + c.BucketURL
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.Object.Get(ctx, cp, nil)
	if err != nil {
		return
	}

	return resp.Body, nil
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	opt := &cos.ObjectGetOptions{
		Range: fmt.Sprintf("bytes=%d-%d", offset, offset+size-1),
	}
	resp, err := c.client.Object.Get(ctx, cp, opt)
	if err != nil {
		return
	}

	return resp.Body, nil
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.Object.Head(ctx, cp, nil)
	if err != nil {
		if e, ok := err.(*cos.ErrorResponse); ok {
			// If object not found, we just need to return a nil object.
			if e.Response.StatusCode == 404 {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}

	// Parse last modified.
	lastModified := convert.StringToTimestamp(resp.Header.Get("Last-Modified"), convert.RFC822)

	o = &model.SingleObject{
		Key:          p,
		Size:         resp.ContentLength,
		LastModified: lastModified,
		MD5:          resp.Header.Get("ETag"),
	}
	return
}
