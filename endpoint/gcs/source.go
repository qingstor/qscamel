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
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := path.Join(c.Path, j.Key) + "/"

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
			object := &model.DirectoryObject{
				Key: utils.Relative(next.Prefix, c.Path),
			}

			fn(object)
			continue
		}

		object := &model.SingleObject{
			Key:  utils.Relative(next.Name, c.Path),
			Size: next.Size,
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
