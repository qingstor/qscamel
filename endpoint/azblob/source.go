package azblob

import (
	"context"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := j.Key
	if cp == "/" {
		cp = ""
	}
	err = c.client.(storage.PrefixLister).ListPrefix(cp,
		pairs.WithObjectFunc(func(object *types.Object) {
			o := &model.SingleObject{
				Key:  object.Name,
				Size: object.Size,
			}

			fn(o)
		}))
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
