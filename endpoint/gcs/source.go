package gcs

import (
	"context"
	"io"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

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

	it := c.client.Objects(ctx, &storage.Query{
		Delimiter: "/",
		Prefix:    cp,
	})
	for {
		next, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if next.Prefix != "" {
			object := model.Object{
				Key:   path.Join(p, path.Base(next.Prefix)),
				IsDir: true,
				Size:  0,
			}

			o = append(o, object)
			continue
		}

		object := model.Object{
			Key:   path.Join(p, path.Base(next.Name)),
			IsDir: false,
			Size:  next.Size,
		}

		o = append(o, object)
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	o := c.client.Object(cp)

	return o.NewReader(ctx)
}
