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
		isDirLink := false
		var size int64
		// if obj is a link, try to get target's metadata
		if v.Mode()&os.ModeSymlink != 0 {
			target, err := getTargetByLink("/" + utils.Join(cp, v.Name()))
			if err != nil {
				return err
			}

			if target.IsDir() {
				isDirLink = true
			} else {
				size = target.Size()
			}
		}

		// TODO: we need to check whether this dir is transferred, in case of cycle link
		if v.IsDir() || isDirLink {
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

		// if size not 0 (v is a link), use target size as o.Size
		if size != 0 {
			o.Size = size
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

// getTargetByLink try to eval symlink and stat the target file
func getTargetByLink(path string) (os.FileInfo, error) {
	tarPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}
	return os.Stat(tarPath)
}
