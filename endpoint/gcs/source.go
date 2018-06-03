package gcs

import (
	"context"
	"path"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.Job, fn func(o *model.Object)) (err error) {

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
			return err
		}
		if next.Prefix != "" {
			object := &model.Object{
				Key:   utils.Relative(next.Prefix, c.Path),
				IsDir: true,
				Size:  0,
			}

			fn(object)
			continue
		}

		object := &model.Object{
			Key:   utils.Relative(next.Name, c.Path),
			IsDir: false,
			Size:  next.Size,
		}

		fn(object)
	}

	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return "", constants.ErrEndpointFuncNotImplemented
}

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return false
}
