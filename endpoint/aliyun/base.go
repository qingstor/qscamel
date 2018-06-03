package aliyun

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	r, err = c.client.GetObject(cp)
	if err != nil {
		return
	}

	return
}

// ReadAt implement source.ReadAt
func (c *Client) ReadAt(
	ctx context.Context, p string, start, end int64,
) (b []byte, err error) {
	cp := utils.Join(c.Path, p)

	r, err := c.client.GetObject(cp, oss.Range(start, end))
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

	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         size,
		LastModified: lastModified,
		MD5:          resp.Get("ETag"),
	}
	return
}
