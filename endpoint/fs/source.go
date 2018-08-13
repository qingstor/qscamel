package fs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := filepath.Join(c.AbsPath, j.Key)

	fi, err := os.Open(cp)
	if err != nil {
		return
	}
	list, err := fi.Readdir(-1)
	fi.Close()

	for _, v := range list {
		if v.IsDir() {
			o := &model.DirectoryObject{
				Key: "/" + utils.Join(j.Key, v.Name()),
			}

			fn(o)

			continue
		}
		o := &model.SingleObject{
			Key:  "/" + utils.Join(j.Key, v.Name()),
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
