package upyun

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/upyun/go-sdk/upyun"
	"io"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}

// Readable implement source.Readable
func (c *Client) Readable() bool {
	return true
}

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.Job, rc chan *model.Object) {
	defer close(rc)

	cp := utils.Join(c.Path, j.Path) + "/"

	oc := make(chan *upyun.FileInfo, 100)

	err := c.client.List(&upyun.GetObjectsConfig{
		Path:         cp,
		MaxListLevel: 1,
		ObjectsChan:  oc,
	})
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		rc <- nil
		return
	}

	for v := range oc {
		rc <- &model.Object{
			Key:   utils.Relative(v.Name, c.Path),
			IsDir: v.IsDir,
			Size:  v.Size,
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	r, w := io.Pipe()

	_, err = c.client.Get(&upyun.GetObjectConfig{
		Path:   cp,
		Writer: w,
	})
	if err != nil {
		return
	}
	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
