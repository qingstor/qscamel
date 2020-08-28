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
		// if v is a link, and client not follow link, skip it
		if v.Mode()&os.ModeSymlink != 0 && !c.EnableLinkFollow {
			continue
		}

		target, err := checkLink(v, cp)
		if err != nil {
			return err
		}

		if target.IsDir() {
			o := &model.DirectoryObject{
				Key: "/" + utils.Join(j.Key, v.Name()), // always use current v's name as key
			}

			fn(o)

			continue
		}
		o := &model.SingleObject{
			Key:  "/" + utils.Join(j.Key, v.Name()), // always use current v's name as key
			Size: target.Size(),
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

// checkLink handle a FileInfo at current path and follow link if needed
func checkLink(v os.FileInfo, cp string) (os.FileInfo, error) {
	// if v is not link, return directly
	if v.Mode()&os.ModeSymlink == 0 {
		return v, nil
	}

	// otherwise, follow the link to get the target
	tarPath, err := filepath.EvalSymlinks("/" + utils.Join(cp, v.Name()))
	if err != nil {
		return nil, err
	}
	return os.Stat(tarPath)
}
