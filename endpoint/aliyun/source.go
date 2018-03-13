package aliyun

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/yunify/qscamel/model"
)

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}

// Readable implement source.Readable
func (c *Client) Readable() bool {
	return true
}

// List implement source.List
func (c *Client) List(ctx context.Context, p string) (o []model.Object, err error) {
	o = []model.Object{}

	// Add "/" to list specific prefix.
	cp := path.Join(c.Path, p) + "/"
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	marker := ""
	first := true

	for marker != "" || first {
		resp, err := c.client.ListObjects(
			oss.Delimiter("/"),
			oss.Marker(marker),
			oss.MaxKeys(MaxKeys),
			oss.Prefix(cp),
		)
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Objects {
			object := model.Object{
				Key:   path.Join(p, path.Base(v.Key)),
				IsDir: false,
				Size:  v.Size,
			}

			o = append(o, object)
		}
		for _, v := range resp.CommonPrefixes {
			object := model.Object{
				Key:   path.Join(p, path.Base(v)),
				IsDir: true,
				Size:  0,
			}

			o = append(o, object)
		}

		first = false
		marker = resp.NextMarker
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	r, err = c.client.GetObject(cp)
	if err != nil {
		return
	}

	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
