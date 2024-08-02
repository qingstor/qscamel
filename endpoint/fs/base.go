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
func (c *Client) Read(ctx context.Context, p string, _ bool) (r io.Reader, err error) {
	cp, err := c.Encode(filepath.Join(c.AbsPath, p))
	if err != nil {
		return
	}

	r, err = os.Open(cp)
	if err != nil {
		return
	}
	return
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp, err := c.Encode(filepath.Join(c.AbsPath, p))
	if err != nil {
		return
	}

	f, err := os.Open(cp)
	if err != nil {
		return
	}

	r = io.NewSectionReader(f, offset, size)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string, _ bool) (o *model.SingleObject, err error) {
	cp, err := c.Encode(filepath.Join(c.AbsPath, p))
	if err != nil {
		return
	}

	fi, err := os.Stat(cp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return
	}
	// We will not calculate md5 while stating object.
	o = &model.SingleObject{
		Key:          p,
		Size:         fi.Size(),
		LastModified: fi.ModTime().Unix(),
	}
	return
}
