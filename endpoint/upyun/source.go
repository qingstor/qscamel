package upyun

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/upyun/go-sdk/upyun"

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
func (c *Client) List(ctx context.Context, p string, rc chan *model.Object) (err error) {
	defer close(rc)

	// Add "/" to list specific prefix.
	cp := path.Join(c.Path, p) + "/"
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	oc := make(chan *upyun.FileInfo, 100)

	err = c.client.List(&upyun.GetObjectsConfig{
		Path:         cp,
		MaxListLevel: 1,
		ObjectsChan:  oc,
	})
	if err != nil {
		return
	}

	for obj := range oc {
		rc <- &model.Object{
			Key:   path.Join(p, path.Base(obj.Name)),
			IsDir: obj.IsDir,
			Size:  obj.Size,
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

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

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
