package upyun

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/upyun/go-sdk/upyun"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "upyun:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string, _ bool) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	r, w := io.Pipe()

	_, err = c.client.Get(&upyun.GetObjectConfig{
		Path:   cp,
		Writer: w,
	})
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

	r, w := io.Pipe()

	_, err = c.client.Get(&upyun.GetObjectConfig{
		Path:   cp,
		Writer: w,
		Headers: map[string]string{
			"Range": fmt.Sprintf("bytes=%d-%d", offset, offset+size-1),
		},
	})

	return r, nil
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string, _ bool) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetInfo(cp)
	if err != nil {
		// If not found, upyun sdk will return an error contains "HEAD 404"
		if strings.Contains(err.Error(), "HEAD 404") {
			return nil, nil
		}
		logrus.Errorf("Get %s info failed for %v.", p, err)
		return
	}
	o = &model.SingleObject{
		Key:          p,
		Size:         resp.Size,
		LastModified: resp.Time.Unix(),
		MD5:          resp.ETag,
	}
	return
}
