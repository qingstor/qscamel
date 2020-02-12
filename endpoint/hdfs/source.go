package hdfs

import (
	"context"
	"path/filepath"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := filepath.Join(c.Path, j.Key)

	fi, err := c.client.ReadDir(cp)
	if err != nil {
		return
	}

	for _, v := range fi {
		if v.IsDir() {
			o := &model.DirectoryObject{
				Key: filepath.Join(j.Key, v.Name()),
			}

			fn(o)

			continue
		}
		o := &model.SingleObject{
			Key:  filepath.Join(j.Key, v.Name()),
			Size: v.Size(),
		}

		fn(o)
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
