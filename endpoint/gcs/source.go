package gcs

import (
	"context"
	"io"
	"path"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
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
func (c *Client) List(ctx context.Context, j *model.Job, rc chan *model.Object) {
	defer close(rc)

	cp := path.Join(c.Path, j.Path) + "/"

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
			logrus.Errorf("List objects failed for %v.", err)
			rc <- nil
			return
		}
		if next.Prefix != "" {
			object := &model.Object{
				Key:   utils.Relative(next.Prefix, c.Path),
				IsDir: true,
				Size:  0,
			}

			rc <- object
			continue
		}

		object := &model.Object{
			Key:   utils.Relative(next.Name, c.Path),
			IsDir: false,
			Size:  next.Size,
		}

		rc <- object
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	o := c.client.Object(cp)

	return o.NewReader(ctx)
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
