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
func (c *Client) List(ctx context.Context, j *model.Job, fn func(o *model.Object)) (err error) {
	cp := filepath.Join(c.AbsPath, j.Path)

	fi, err := os.Open(cp)
	if err != nil {
		return
	}
	list, err := fi.Readdir(-1)
	fi.Close()

	for _, v := range list {
		o := &model.Object{
			Key:   "/" + utils.Join(j.Path, v.Name()),
			IsDir: v.IsDir(),
			Size:  v.Size(),
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
