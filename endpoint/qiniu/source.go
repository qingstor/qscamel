package qiniu

import (
	"context"
	"io"
	"path"
	"strings"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
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
func (c *Client) List(ctx context.Context, p string, rc chan *model.Object) (err error) {
	defer close(rc)

	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	// Add "/" to list specific prefix.
	cp := path.Join(c.Path, p) + "/"
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	marker := t.Marker
	first := true

	for marker != "" || first {
		entries, _, nextMarker, _, err := c.bucket.ListFiles(c.BucketName, cp, "", marker, MaxListFileLimit)
		if err != nil {
			logrus.Errorf("List files failed for %v.", err)
			return err
		}
		for _, v := range entries {
			object := &model.Object{
				Key:   strings.TrimLeft(v.Key, c.Path),
				IsDir: false,
				Size:  v.Fsize,
			}

			rc <- object
		}

		first = false
		marker = nextMarker

		// Update task content.
		t.Marker = marker
		err = t.Save(ctx)
		if err != nil {
			logrus.Errorf("Save task failed for %v.", err)
			return err
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

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
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	deadline := time.Now().Add(time.Hour).Unix()
	url = storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)
	return
}
