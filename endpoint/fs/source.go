package fs

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
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
func (c *Client) List(ctx context.Context, j *model.Job, rc chan *model.Object) {
	defer close(rc)

	cp := utils.Join(c.Path, j.Path)

	fi, err := os.Open(cp)
	if err != nil {
		logrus.Errorf("Open dir failed for %v.", err)
		rc <- nil
		return
	}
	list, err := fi.Readdir(-1)
	fi.Close()

	for _, v := range list {
		rc <- &model.Object{
			Key:   utils.Relative(v.Name(), c.Path),
			IsDir: v.IsDir(),
			Size:  v.Size(),
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

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
