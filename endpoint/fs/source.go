package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return false
}

// Readable implement source.Readable
func (c *Client) Readable() bool {
	return true
}

// List implement source.List
func (c *Client) List(ctx context.Context, p string) (o []model.Object, err error) {
	cp := path.Join(c.Path, p)

	fi, err := os.Open(cp)
	if err != nil {
		return nil, err
	}
	list, err := fi.Readdir(-1)
	fi.Close()

	o = make([]model.Object, len(list))
	for k, v := range list {
		o[k] = model.Object{
			Key:   path.Join(p, v.Name()),
			IsDir: v.IsDir(),
			Size:  v.Size(),
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)

	r, err = os.Open(cp)
	if err != nil {
		logrus.Errorf("Fs open file %s failed for %s.", cp, err)
		return
	}
	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
