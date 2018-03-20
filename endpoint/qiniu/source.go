package qiniu

import (
	"context"
	"io"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/sirupsen/logrus"

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

	marker := j.Marker

	for {
		entries, _, nextMarker, _, err := c.bucket.ListFiles(c.BucketName, cp, "", marker, MaxListFileLimit)
		if err != nil {
			logrus.Errorf("List files failed for %v.", err)
			rc <- nil
			return
		}
		for _, v := range entries {
			object := &model.Object{
				Key:   utils.Relative(v.Key, c.Path),
				IsDir: false,
				Size:  v.Fsize,
			}

			rc <- object
		}

		marker = nextMarker

		// Update task content.
		j.Marker = marker
		err = j.Save(ctx)
		if err != nil {
			logrus.Errorf("Save task failed for %v.", err)
			rc <- nil
			return
		}

		if marker == "" {
			break
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)

	resp, err := c.client.Get(url)
	if err != nil {
		return
	}

	r = resp.Body
	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	cp := utils.Join(c.Path, p)

	deadline := time.Now().Add(time.Hour).Unix()
	url = storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)
	return
}
