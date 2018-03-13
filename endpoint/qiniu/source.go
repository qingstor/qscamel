package qiniu

import (
	"context"
	"io"
	"path"
	"strings"
	"time"

	"github.com/qiniu/api.v7/storage"

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
		entries, commonPrefix, nextMarker, _, err := c.bucket.ListFiles(c.BucketName, cp, "/", marker, MaxListFileLimit)
		if err != nil {
			return nil, err
		}
		for _, v := range entries {
			object := model.Object{
				Key:   path.Join(p, path.Base(v.Key)),
				IsDir: false,
				Size:  v.Fsize,
			}

			o = append(o, object)
		}
		for _, v := range commonPrefix {
			object := model.Object{
				Key:   path.Join(p, path.Base(v)),
				IsDir: true,
				Size:  0,
			}

			o = append(o, object)
		}

		first = false
		marker = nextMarker
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)

	resp, err := c.client.Get(url)
	if err != nil {
		return
	}

	r = resp.Body
	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	deadline := time.Now().Add(time.Hour).Unix()
	url = storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)
	return
}
