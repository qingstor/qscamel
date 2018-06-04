package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/yunify/qscamel/model"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "fs:" + c.AbsPath
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := filepath.Join(c.AbsPath, p)

	r, err = os.Open(cp)
	if err != nil {
		return
	}
	return
}

// ReadAt implement source.ReadAt
func (c *Client) ReadAt(
	ctx context.Context, p string, start, end int64,
) (b []byte, err error) {
	cp := filepath.Join(c.AbsPath, p)

	r, err := os.Open(cp)
	if err != nil {
		return
	}
	defer r.Close()

	b = make([]byte, end-start+1)
	_, err = r.ReadAt(b, start)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := filepath.Join(c.AbsPath, p)

	fi, err := os.Stat(cp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return
	}
	// We will not calculate md5 while stating object.
	o = &model.Object{
		Key:          p,
		IsDir:        fi.IsDir(),
		Size:         fi.Size(),
		LastModified: fi.ModTime().Unix(),
	}
	return
}
